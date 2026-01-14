package services

import (
	"CardFlow/internal/config"
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"CardFlow/internal/utils"
	"context"
	"errors"
	"time"
)

type CardService interface {
    CreateCard(context.Context, models.CreateCardReq)(any, error)
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