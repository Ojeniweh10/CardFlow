package services

import (
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"CardFlow/internal/utils"
	"context"
	"errors"
)

type KycService interface {
    Uploadimage(context.Context, models.KycProfile)error
    UploadKycDocument(context.Context, models.KycDoc) error
	UploadProofOfAddress(context.Context, models.KycDoc) error
}

type kycService struct {
    userRepo repositories.UserRepository
    kycrepo repositories.KycRepository
}

func NewKycService(kycrepo repositories.KycRepository, userRepo repositories.UserRepository) KycService {
    return &kycService{kycrepo:kycrepo, userRepo: userRepo}
}

const (
    DocTypeSelfie        = "selfie"
    DocTypeIDDocument    = "id_document"
    DocTypeProofOfAddr   = "proof_of_address"
	DocsUploaded 		 = "documents_uploaded"
	Verified		 = "verified"
)

func (s *kycService) Uploadimage(ctx context.Context,data models.KycProfile) error {
	return s.kycrepo.RunInTransaction(ctx, func(repo repositories.KycRepository) error {

		existing, err := repo.FindByUserID(data.Userid)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
		if existing != nil {
			return errors.New("kyc submission already exists")
		}

		encrypted, mime, err := utils.EncryptBase64Document(data.ImageStr)
		if err != nil {
			return err
		}

		sub := &models.KYCSubmission{
			UserID: data.Userid,
			Status: "started",
		}
		if err := repo.CreateKycSubmission(sub); err != nil {
			return err
		}

		doc := &models.KYCDocument{
			KYCSubmissionID: sub.ID,
			DocumentType:    DocTypeSelfie,
			MimeType:        mime,
			EncryptedData:   []byte(encrypted),
		}

		return repo.CreateKycDocsSubmission(doc)
	})
}


func (s *kycService) UploadKycDocument(ctx context.Context, data models.KycDoc) error {
	return s.kycrepo.RunInTransaction(ctx, func(repo repositories.KycRepository) error {

		sub, err := repo.FindByUserID(data.Userid)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
		if sub == nil {
			return errors.New("kyc submission does not exist, upload selfie first")
		}

		encrypted, mime, err := utils.EncryptBase64Document(data.DocStr)
		if err != nil {
			return err
		}

		sub.Status = DocsUploaded
		if err := repo.UpdateKycSubmission(sub); err != nil {
			return err
		}

		doc := &models.KYCDocument{
			KYCSubmissionID: sub.ID,
			DocumentType:    DocTypeIDDocument,
			MimeType:        mime,
			EncryptedData:   []byte(encrypted),
		}

		return repo.CreateKycDocsSubmission(doc)
	})
}


func (s *kycService) UploadProofOfAddress(ctx context.Context, data models.KycDoc) error {
	return s.kycrepo.RunInTransaction(ctx, func(repo repositories.KycRepository) error {

		sub, err := repo.FindByUserID(data.Userid)
		if err != nil {
			return errors.New("something went wrong, please try again later")
		}
		if sub == nil {
			return errors.New("kyc submission does not exist, upload selfie first")
		}

		encrypted, mime, err := utils.EncryptBase64Document(data.DocStr)
		if err != nil {
			return err
		}

		sub.Status = Verified
		if err := repo.UpdateKycSubmission(sub); err != nil {
			return err
		}

		doc := &models.KYCDocument{
			KYCSubmissionID: sub.ID,
			DocumentType:    DocTypeProofOfAddr,
			MimeType:        mime,
			EncryptedData:   []byte(encrypted),
		}

		return repo.CreateKycDocsSubmission(doc)
	})
}


//make sure you fully understand the code before going ahead with other coding.