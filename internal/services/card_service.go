package services

import (
	"CardFlow/internal/config"
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"CardFlow/internal/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CardService interface {
    CreateCard(context.Context, models.CreateCardReq)(any, error)
	GetAllCards(context.Context, uuid.UUID)([]models.GetAllCardsResp, error)
	GetCardById(context.Context, models.GetCardReq)(models.GetCardResp, error)
	ModifyCardStatus(ctx context.Context, data models.GetCardReq, status string) error
	TopUpCard(ctx context.Context, data models.TopUpCardReq)(any, error)
}

type cardService struct {
	userrepo repositories.UserRepository
    kycrepo repositories.KycRepository
	cardrepo repositories.CardRepository
	Txnrepo repositories.TransactionRepository
}

func NewCardService(userRepo repositories.UserRepository,  kycrepo repositories.KycRepository, cardRepo repositories.CardRepository, txnRepo repositories.TransactionRepository) CardService {
    return &cardService{userrepo:userRepo, kycrepo:kycrepo, cardrepo: cardRepo, Txnrepo:txnRepo}
}

var ErrUserNotFound = errors.New("user not found")


func (s *cardService) CreateCard(ctx context.Context, data models.CreateCardReq)(any, error){
	var ExpiryMonth, ExpiryYear string
	var ExpiresAt time.Time
	//first check if kyc for user has been verified
	user, err := s.kycrepo.FindByUserID(data.Userid)
	if err != nil{
		return nil, errors.New("Something Went Wrong, Please try again later")
	}
	if user.Status != "verified"{
		return nil, errors.New("incomplete kyc")
	}
	//set expiry month and year to 1 year from now if its single use or 3 years if multi use
	switch data.CardType {
	case "single-use":
		ExpiryMonth, ExpiryYear, ExpiresAt = utils.GetExpiryDate(1)
	case "multi-use":
		ExpiryMonth, ExpiryYear, ExpiresAt = utils.GetExpiryDate(3)
	default:
		return nil, errors.New("invalid card type")
	}
	//generate card reference
	CardReference := GenerateCardReference("CRDFLW")
	//generate card
	IIN := config.IIN
	MiddleDigits := utils.GenerateNumberString(8)
	CardNumber := "4" + IIN + MiddleDigits
	Cvv := utils.GenerateNumberString(3)
	MaskedCardNumber := "4" + IIN + MiddleDigits[:4] + "****" + MiddleDigits[6:]
	//luhn check digit
	CheckDigit := utils.ComputeLuhnCheckDigit(CardNumber)
	FullCardNumber := CardNumber + CheckDigit
	//encrypt card number
	PANENcrypted, err := utils.EncryptString(FullCardNumber, config.EncryptionKey)
	if err != nil{
		return nil, errors.New("Something Went Wrong, Please try again later")
	}
	CvvEncrypted, err := utils.EncryptString(Cvv, config.EncryptionKey)
	if err != nil{
		return nil, errors.New("Something Went Wrong, Please try again later")
	}
	//hash card number
	HashedCardNumber, err := utils.HashString(FullCardNumber)
	if err != nil{
		return nil, errors.New("Something Went Wrong, Please try again later")
	}
	HashedCvv, err := utils.HashString(Cvv)
	if err != nil{
		return nil, errors.New("Something Went Wrong, Please try again later")
	}
	card := &models.Card{
		UserID:        data.Userid,
		CardType:     data.CardType,
		Currency:     data.Currency,
		SpendingLimitAmount: data.SpendingLimit,
		CardReference: CardReference,
		PANencrypted: PANENcrypted,
		CVVencrypted: CvvEncrypted,
		PANHash: HashedCardNumber,
		CVVHash: HashedCvv,
		MaskedPAN: MaskedCardNumber,
		LastFour: FullCardNumber[len(FullCardNumber)-4:],
		Status: "active",
		ExpiryMonth: ExpiryMonth,
		ExpiryYear: ExpiryYear,
		ExpiresAt: ExpiresAt,
	}
	err = s.cardrepo.CreateCard(ctx, card)
	if err != nil{
		return nil, errors.New("Something Went Wrong, Please try again later")
	}
	
	resp := &models.CreateCardResp{
		CardType: data.CardType,
		MaskedPAN: MaskedCardNumber,
		Currency: data.Currency,
		SpendingLimit: data.SpendingLimit,
		Balance: 0.00,
		Cvv: Cvv,
		Status: "active",
		ExpiryMonth: ExpiryMonth,
		ExpiryYear: ExpiryYear,
	}

	return resp, nil
}

func GenerateCardReference(prefix string) string {
	uniqueID := utils.GenerateRandomString(10)
	return prefix + uniqueID
}

func (s *cardService) GetAllCards(ctx context.Context, Userid uuid.UUID)([]models.GetAllCardsResp, error){
	var res []models.GetAllCardsResp
	cards, err := s.cardrepo.FindCardsByID(ctx, Userid)
	if err != nil {
		return nil,  errors.New("something went wrong, please try again later")
	}
	for _, card := range cards{
		CvvDecrypt, err := utils.DecryptString(card.CVVencrypted, config.EncryptionKey)
		if err != nil {
			return nil, errors.New("Something Went Wrong, Please try again later")
		}
		resp := models.GetAllCardsResp{
			Cardid: card.ID,
			CardType: card.CardType,
			MaskedPAN: card.MaskedPAN,
			Lastfour: card.LastFour,
			Cvv: CvvDecrypt,
			Currency: card.Currency,
			Status: card.Status,
		}
		res = append(res, resp)
	}
	
	return res, nil
}

func (s *cardService) GetCardById(ctx context.Context, data models.GetCardReq)(models.GetCardResp, error){
	card, err := s.cardrepo.FindCardByID(ctx, data)
	if err != nil {
		return models.GetCardResp{},  errors.New("something went wrong, please try again later")
	}
	Pan, err := utils.DecryptString(card.PANencrypted, config.EncryptionKey)
	if err != nil {
		return models.GetCardResp{}, errors.New("Something Went Wrong, Please try again later")
	}
	Cvv, err := utils.DecryptString(card.CVVencrypted, config.EncryptionKey)
	if err != nil {
		return models.GetCardResp{}, errors.New("Something Went Wrong, Please try again later")
	}
	res := models.GetCardResp{
		Cardid: card.ID,
		CardType: card.CardType,
		PAN: Pan,
		Cvv: Cvv,
		Lastfour: card.LastFour,
		Currency: card.Currency,
		Status: card.Status,
		SpendingLimit: card.SpendingLimitAmount,
		CurrentBalance: card.CurrentBalance,
		ExpiryMonth: card.ExpiryMonth,
		ExpiryYear: card.ExpiryYear,
	}
	
	return res, nil
}

func (s *cardService) ModifyCardStatus(ctx context.Context, data models.GetCardReq, status string) error{
	switch status {
	case "freeze":
		card, err := s.cardrepo.FindCardByID(ctx, data)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
		switch card.Status {
		case "frozen":
			return errors.New("card is already frozen")
		case "expired":
			return errors.New("card has expired")
		}
		card.Status = "frozen"
		err = s.cardrepo.Update(ctx, card)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
	case "unfreeze":
		card, err := s.cardrepo.FindCardByID(ctx, data)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
		switch card.Status {
		case "active":
			return errors.New("card is already active")
		case "expired":
			return errors.New("card has expired")
		}
		card.Status = "active"
		err = s.cardrepo.Update(ctx, card)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
	case "terminate":
		card, err := s.cardrepo.FindCardByID(ctx, data)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
		switch card.Status {
		case "terminated":
			return errors.New("card is already terminated")
		case "expired":
			return errors.New("card has expired")
		}
		card.Status = "terminated"
		err = s.cardrepo.Update(ctx, card)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
	}
	
	return nil
}

func (s *cardService) TopUpCard(ctx context.Context, data models.TopUpCardReq)(any, error){
	cardid, err := uuid.Parse(data.Cardid)
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	user, err := s.userrepo.FindByID(ctx, data.Userid)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("something went wrong")
	}
	cardreq := &models.GetCardReq{
		UserId: data.Userid,
		CardId: data.Cardid,
	}
	card , err := s.cardrepo.FindCardByID(ctx, *cardreq)
	if err != nil {
		//log err to devs
		return nil, errors.New("something went wrong")
	}
	//figure out how to check if card variable that stores struct of models.card is empty
	if card.ID == uuid.Nil {
		return nil, errors.New("card not found")
	}
	switch card.Status{
		case "frozen":
			return nil,  errors.New("card is already frozen")
		case "expired":
			return nil, errors.New("card has expired")
		case "terminated":
			return nil, errors.New("card is already terminated")		
	}
	fee := data.Amount * 0.01
	newAmount := data.Amount - fee 
	card.CurrentBalance = card.CurrentBalance + newAmount
	err = s.cardrepo.Update(ctx, card)
	if err != nil {
		return nil,  errors.New("something went wrong, please try again later")
	}
	transaction_reference := GenerateCardReference("tOP-UP")
	transactions := &models.Transaction{
		UserID: card.UserID,
		CardID: card.ID,
		TransactionReference: transaction_reference,
		Amount: data.Amount,
		Currency: "USD",
		Type: "funding",
		Direction: "credit",
		Status: "completed",
		TransactionTimestamp: time.Now(),
	}
	err = s.Txnrepo.CreateTransaction(ctx, transactions)
	if err != nil{
		return nil, errors.New("something went wrong, please try again")
	}
	Balanceledger := &models.BalanceLedger{
		CardID: cardid,
		TransactionID: transactions.ID,
		EntryType: "card top-up",
		Amount: data.Amount,
		FeeCharged : fee,
		BalanceAfter: card.CurrentBalance,
	}
	err = s.Txnrepo.CreateLedger(ctx, *Balanceledger)
	if err != nil {
		return nil,  errors.New("something went wrong, please try again later")
	}
	//notify user via email card has been funded
	res := map[string]string{
		"firstname": user.FirstName,
		"email": user.Email,
		"lastfour": card.LastFour,
		"amount": fmt.Sprintf("%f",data.Amount),
		"fee": fmt.Sprintf("%f",fee),
	}
	err = SendCardTopUpEmail(res)
	if err != nil {
		//log error to devs and retry automatically later via a queue to send the email.
	}
	
	return nil, nil
}
