package models

import (
	"time"

	"github.com/google/uuid"
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