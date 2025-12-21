package services

import (
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"CardFlow/internal/utils"
	"errors"
)

type KycService interface {
    Uploadimage(models.KycProfile)error
    VerifyBVN(models.Bvn) error
}

type kycService struct {
    repo repositories.KycRepository
}

func NewKycService(repo repositories.KycRepository) KycService {
    return &kycService{repo:repo}
}

func (s *kycService)Uploadimage(data models.KycProfile) error{
    kycProfile, err := s.repo.FindByID(data.Userid)
	if kycProfile != nil {
		return errors.New("kyc profile already exists for this user")
	}
	if err != nil {
		//log the error to notify devs then return a generic error message
		return errors.New("something went wrong, please try again later")
	}
    userid := data.Userid.String()
    imageUrl, err := utils.ProcessBase64File(data.ImageStr, "kyc image", userid)
	if err != nil {
		return err
	}
    userKyc := &models.KYCProfile{
        UserID: data.Userid,
        DateOfBirth: data.DOB,
        ImageURL: imageUrl,
        Status: "started",
    }

    return s.repo.CreateKycProfile(userKyc)

}

func (s *kycService)VerifyBVN(data models.Bvn)error{
    //take bvn, call partner, take results, run matches on system if it doesnt match don't save and return reason for rejeciton else save.
    return  nil
}