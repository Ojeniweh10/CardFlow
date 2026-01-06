package services

import (
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"context"
	"testing"

	"github.com/google/uuid"
)

type fakeKycRepo struct {
	existingSubmission *models.KYCSubmission
	findErr            error
	createErr          error
	updateErr          error
	createDocErr       error
}

func (f *fakeKycRepo) RunInTransaction(ctx context.Context, fn func(repo repositories.KycRepository) error) error {
	return fn(f) // just call the callback with itself
}

func (f *fakeKycRepo) FindByUserID(userID uuid.UUID) (*models.KYCSubmission, error) {
	return f.existingSubmission, f.findErr
}

func (f *fakeKycRepo) CreateKycSubmission(sub *models.KYCSubmission) error {
	return f.createErr
}

func (f *fakeKycRepo) UpdateKycSubmission(sub *models.KYCSubmission) error {
	return f.updateErr
}

func (f *fakeKycRepo) CreateKycDocsSubmission(doc *models.KYCDocument) error {
	return f.createDocErr
}


func TestUploadImage_Success(t *testing.T) {
	repo := &fakeKycRepo{
		existingSubmission: nil, // no submission exists yet
	}

	service := &kycService{kycrepo: repo}

	req := models.KycProfile{
		Userid:   uuid.New(),
		ImageStr: "fake-base64-image",
	}

	err := service.Uploadimage(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUploadImage_AlreadyExists(t *testing.T) {
	repo := &fakeKycRepo{
		existingSubmission: &models.KYCSubmission{UserID: uuid.New()},
	}

	service := &kycService{kycrepo: repo}

	req := models.KycProfile{
		Userid:   uuid.New(),
		ImageStr: "fake-base64-image",
	}

	err := service.Uploadimage(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error because submission exists, got nil")
	}
}


