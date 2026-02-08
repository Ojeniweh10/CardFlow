package services

import (
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"CardFlow/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)


type TransactionService interface{
	WebhookTransaction(ctx context.Context, data models.WebhookReq)(any, error)
}

type transactionService struct {
	userrepo repositories.UserRepository
	cardrepo repositories.CardRepository
	Txnrepo repositories.TransactionRepository
}
func NewTransactionService(Txnrepo repositories.TransactionRepository, cardRepo repositories.CardRepository, userRepo repositories.UserRepository) TransactionService {
    return &transactionService{Txnrepo:Txnrepo, cardrepo: cardRepo, userrepo: userRepo}
}


func (s *transactionService) WebhookTransaction(ctx context.Context,data models.WebhookReq) (any, error) {

	// --------------------------------------------------
	// 1. Idempotency Guard
	// --------------------------------------------------
	if data.IdempotencyKey != "" {
		existing, _ := s.Txnrepo.FindByIdempotencyKey(ctx, data.IdempotencyKey)
		if existing.ID != uuid.Nil {
			// Webhook already processed — acknowledge safely
			return map[string]string{"status": "duplicate_ignored"}, nil
		}
	}

	// --------------------------------------------------
	// 2. Load Card
	// --------------------------------------------------
	card, err := s.cardrepo.FindCardsByReference(ctx, data)
	if err != nil || card.ID == uuid.Nil {
		return nil, errors.New("card not found")
	}
	switch card.Status {
	case "frozen", "expired", "terminated":
		return nil, errors.New("card is not active")
	}

	availableBalance := card.CurrentBalance - card.HeldBalance

	switch data.Type {
	case "authorization":

		if availableBalance < data.Amount {
			return nil, errors.New("insufficient available balance")
		}

		if card.SpendingLimitAmount < data.Amount {
			return nil, errors.New("exceeds card spending limit")
		}

		// Create transaction
		txn := &models.Transaction{
			UserID:               card.UserID,
			CardID:               card.ID,
			TransactionReference: data.TransactionID,
			IdempotencyKey:       &data.IdempotencyKey,
			Amount:               data.Amount,
			Currency:             data.Currency,
			AuthorizedAmount:     data.Amount,
			CapturedAmount:       0,
			Type:                 "authorization",
			Direction:            "debit",
			Status:               "authorized",
			MerchantName:         &data.Merchant.Name,
			MerchantMCC:          &data.Merchant.MCC,
			MerchantCountry:      &data.Merchant.Country,
			Source:               &data.Network,
			TransactionTimestamp: data.Timestamp,
		}

		if err := s.Txnrepo.CreateTransaction(ctx, txn); err != nil {
			return nil, err
		}

		// Increase held balance
		card.HeldBalance += data.Amount
		if err := s.cardrepo.Update(ctx, card); err != nil {
			return nil, err
		}

		// Ledger entry (HOLD)
		ledger := &models.BalanceLedger{
			CardID:        card.ID,
			TransactionID: txn.ID,
			EntryType:     "Authorization Hold",
			Amount:        data.Amount,
			FeeCharged:    0,
			BalanceAfter:  card.CurrentBalance,
		}
		_ = s.Txnrepo.CreateLedger(ctx, *ledger)

		return map[string]string{"status": "authorized"}, nil

	case "capture":
		user, err := s.userrepo.FindByID(ctx, card.UserID)
		if err != nil {
			if errors.Is(err, ErrUserNotFound) {
				return nil, errors.New("user not found")
			}
			return nil, errors.New("something went wrong")
		}
		txn, err := s.Txnrepo.FindTxnByReference(ctx, data.TransactionID)
		if err != nil || txn.ID == uuid.Nil {
			return nil, errors.New("transaction not found")
		}

		if txn.Status != "authorized" {
			return nil, errors.New("transaction not eligible for capture")
		}

		fee := data.Amount * 0.01

		// Release hold & deduct balance
		card.HeldBalance -= txn.AuthorizedAmount
		card.CurrentBalance -= (data.Amount + fee)

		if err := s.cardrepo.Update(ctx, card); err != nil {
			return nil, err
		}

		// Update transaction
		txn.CapturedAmount = data.Amount
		txn.Status = "completed"
		txn.Type = "capture"
		txn.TransactionTimestamp = data.Timestamp

		if err := s.Txnrepo.Update(ctx, txn); err != nil {
			return nil, err
		}

		ledger := &models.BalanceLedger{
			CardID:        card.ID,
			TransactionID: txn.ID,
			EntryType:     "Capture Settlement",
			Amount:        data.Amount,
			FeeCharged:    fee,
			BalanceAfter:  card.CurrentBalance,
		}
		_ = s.Txnrepo.CreateLedger(ctx, *ledger)

		// Notify user
		res := map[string]string{
			"firstname": user.FirstName,
			"email":     user.Email,
			"lastfour":  card.LastFour,
			"amount":    fmt.Sprintf("%.2f", data.Amount),
			"fee":       fmt.Sprintf("%.2f", fee),
			"balance":   fmt.Sprintf("%.2f", card.CurrentBalance),
		}
		err = utils.SendWithRetry(3, 2*time.Second, func() error {
		return SendCardDebitEmail(res)
		})
		if err != nil {
			// DO NOT return error
			// Just log for observability
			log.Printf(
				"failed to send top-up email for card %s: %v",
				card.ID,
				err,
			)
		}
		return map[string]string{"status": "captured"}, nil

	case "reversal":

		txn, err := s.Txnrepo.FindTxnByReference(ctx, data.TransactionID)
		if err != nil || txn.ID == uuid.Nil {
			return nil, errors.New("transaction not found")
		}

		if txn.Status != "authorized" {
			return nil, errors.New("only authorized transactions can be reversed")
		}

		// Release hold
		card.HeldBalance -= txn.AuthorizedAmount
		if err := s.cardrepo.Update(ctx, card); err != nil {
			return nil, err
		}

		txn.Status = "reversed"
		txn.Type = "reversal"
		txn.TransactionTimestamp = data.Timestamp
		_ = s.Txnrepo.Update(ctx, txn)

		ledger := &models.BalanceLedger{
			CardID:        card.ID,
			TransactionID: txn.ID,
			EntryType:     "Authorization Reversal",
			Amount:        txn.AuthorizedAmount,
			FeeCharged:    0,
			BalanceAfter:  card.CurrentBalance,
		}
		_ = s.Txnrepo.CreateLedger(ctx, *ledger)

		return map[string]string{"status": "reversed"}, nil

	case "refund":

		origTxn, err := s.Txnrepo.FindTxnByReference(ctx, data.OriginalTransactionID)
		if err != nil || origTxn.ID == uuid.Nil {
			return nil, errors.New("original transaction not found")
		}

		if origTxn.Status != "completed" {
			return nil, errors.New("only completed transactions can be refunded")
		}
		user, err := s.userrepo.FindByID(ctx, card.UserID)
		if err != nil || user.ID == uuid.Nil {
			return nil, errors.New("card user not found")
		}

		refundTxn := &models.Transaction{
			UserID:               card.UserID,
			CardID:               card.ID,
			TransactionReference: data.TransactionID,
			IdempotencyKey:       &data.IdempotencyKey,
			Amount:               data.Amount,
			Currency:             data.Currency,
			Type:                 "refund",
			Direction:            "credit",
			Status:               "completed",
			TransactionTimestamp: data.Timestamp,
		}
		_ = s.Txnrepo.CreateTransaction(ctx, refundTxn)

		card.CurrentBalance += data.Amount
		_ = s.cardrepo.Update(ctx, card)

		ledger := &models.BalanceLedger{
			CardID:        card.ID,
			TransactionID: refundTxn.ID,
			EntryType:     "Refund",
			Amount:        data.Amount,
			FeeCharged:    0,
			BalanceAfter:  card.CurrentBalance,
		}
		_ = s.Txnrepo.CreateLedger(ctx, *ledger)

		res := map[string]string{
			"firstname": user.FirstName,
			"email":    user.Email,
			"lastfour": card.LastFour,
			"amount":   fmt.Sprintf("%.2f", data.Amount),
			"balance":  fmt.Sprintf("%.2f", card.CurrentBalance),
		}
		err = utils.SendWithRetry(3, 2*time.Second, func() error {
		return SendRefundEmail(res)
		})
		if err != nil {
			// DO NOT return error
			// Just log for observability
			log.Printf(
				"failed to send top-up email for card %s: %v",
				card.ID,
				err,
			)
		}

		return map[string]string{"status": "refunded"}, nil
	}

	return nil, errors.New("unsupported webhook type")
}