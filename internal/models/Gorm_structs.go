package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)


type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"column:password_hash;size:255;not null"`

	FirstName string `gorm:"column:first_name;size:100;not null"`
	LastName  string `gorm:"column:last_name;size:100;not null"`

	Phone string `gorm:"size:20"`

	Status string `gorm:"size:50;default:active"`

	EmailVerified bool `gorm:"default:false"`
	OTP           string `gorm:"column:otp;size:6"`
	OTPExpiresAt  time.Time
	MFAEnabled    bool `gorm:"column:mfa_enabled;default:false"`
	MFASecret     string `gorm:"column:mfa_secret;size:255"`

	LastLoginAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Admin struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Email		string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"column:password_hash;size:255;not null"`
	
	FirstName string `gorm:"column:first_name;size:100;not null"`
	LastName  string `gorm:"column:last_name;size:100;not null"`

	Phone *string `gorm:"size:20"`

	Role string `gorm:"size:50;default:admin"`
	Status string `gorm:"size:50;default:active"`

	EmailVerified bool `gorm:"default:false"`
	MFAEnabled    bool `gorm:"column:mfa_enabled;default:false"`
	MFASecret     *string `gorm:"column:mfa_secret;size:255"`

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type KYCProfile struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	DateOfBirth string `gorm:"column:date_of_birth;type:text;not null"`
	ImageURL    string    `gorm:"column:image_url;type:text;not null"`

	Status string `gorm:"size:50;default:started"`

	CreatedAt time.Time
	UpdatedAt time.Time
}


type KYCVerification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	KYCProfileID uuid.UUID `gorm:"type:uuid;not null;index"`

	Type string `gorm:"size:20;not null"` // 'bvn' or 'nin'

	Identifier string `gorm:"size:100;not null"` // BVN or NIN

	Status string `gorm:"size:50;default:pending"`

	PartnerRequestID string `gorm:"size:255"`
	PartnerReference string `gorm:"size:255"`

	NameMatch bool
	FaceMatch bool

	PartnerAddressJSON     datatypes.JSON `gorm:"column:partner_address_json;type:jsonb"`
	VerificationResultJSON datatypes.JSON `gorm:"column:verification_result_json;type:jsonb"`

	RejectionReason string `gorm:"type:text"`

	SubmittedAt time.Time `gorm:"default:current_timestamp"`
	VerifiedAt  *time.Time
	ExpiresAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
