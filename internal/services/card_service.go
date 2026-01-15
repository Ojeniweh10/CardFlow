package services

import (
	"CardFlow/internal/config"
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"CardFlow/internal/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type CardService interface {
    CreateCard(context.Context, models.CreateCardReq)(any, error)
	GetAllCards(context.Context, uuid.UUID)([]models.GetAllCardsResp, error)
	GetCardById(context.Context, models.GetCardReq)(models.GetCardResp, error)
}

type cardService struct {
    kycrepo repositories.KycRepository
	cardrepo repositories.CardRepository
}

func NewCardService(kycrepo repositories.KycRepository, cardRepo repositories.CardRepository) CardService {
    return &cardService{kycrepo:kycrepo, cardrepo: cardRepo}
}

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