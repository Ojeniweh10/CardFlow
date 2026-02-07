package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

//
// =========================
// Users
// =========================
//

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"column:password_hash;size:255;not null"`

	FirstName string `gorm:"column:first_name;size:100;not null"`
	LastName  string `gorm:"column:last_name;size:100;not null"`

	Phone string `gorm:"size:20"`

	Status string `gorm:"size:50;not null;default:active"`

	EmailVerified bool `gorm:"not null;default:false"`
	OTP           string    `gorm:"column:otp;size:6"`
	OTPExpiresAt  time.Time `gorm:"column:otp_expires_at"`

	MFAEnabled bool    `gorm:"column:mfa_enabled;not null;default:false"`
	MFASecret  string `gorm:"column:mfa_secret;size:255"`

	LastLoginAt *time.Time `gorm:"column:last_login_at"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

//
// =========================
// Admins
// =========================
//

type Admin struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"column:password_hash;size:255;not null"`

	FirstName string `gorm:"column:first_name;size:100;not null"`
	LastName  string `gorm:"column:last_name;size:100;not null"`

	Phone *string `gorm:"size:20"`

	Role   string `gorm:"size:50;not null;default:admin"`
	Status string `gorm:"size:50;not null;default:active"`

	EmailVerified bool    `gorm:"not null;default:false"`
	MFAEnabled    bool    `gorm:"column:mfa_enabled;not null;default:false"`
	MFASecret     *string `gorm:"column:mfa_secret;size:255"`

	LastLoginAt *time.Time `gorm:"column:last_login_at"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

//
// =========================
// KYC
// =========================
//

type KYCSubmission struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	Status           string  `gorm:"size:50;not null"`
	RejectionReason  *string `gorm:"type:text"`

	ReviewedBy *uuid.UUID `gorm:"type:uuid"`
	Reviewer   *Admin    `gorm:"foreignKey:ReviewedBy"`

	SubmittedAt time.Time  `gorm:"not null;default:current_timestamp"`
	ReviewedAt  *time.Time
	ExpiresAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type KYCDocument struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	KYCSubmissionID uuid.UUID      `gorm:"type:uuid;not null;index, unique"`
	KYCSubmission   KYCSubmission  `gorm:"foreignKey:KYCSubmissionID"`

	DocumentType      string `gorm:"size:50;not null, unique"`
	MimeType          string `gorm:"size:100;not null"`
	EncryptedData     []byte `gorm:"type:bytea;not null"`
	EncryptionVersion string `gorm:"size:20;not null;default:v1"`

	CreatedAt time.Time
}

//
// =========================
// Cards
// =========================
//

type Card struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	CardReference string `gorm:"size:100;uniqueIndex;not null"`
	PANencrypted  string `gorm:"column:pan_encrypted;size:255;not null"`
	CVVencrypted  string `gorm:"column:cvv_encrypted;size:255;not null"`
	MaskedPAN     string `gorm:"column:masked_pan;size:255;not null"`
	LastFour      string `gorm:"column:last_four;size:4;not null"`
	PANHash       string `gorm:"column:pan_hash;size:255;not null"`
	CVVHash       string `gorm:"column:cvv_hash;size:255;not null"`

	CardType string `gorm:"size:50"`
	Currency string `gorm:"size:3;not null;default:USD"`

	Status string `gorm:"size:50"`

	SpendingLimitAmount float64 `gorm:"type:decimal(15,2)"`

	CurrentBalance float64 `gorm:"type:decimal(15,2);not null;default:0.00"`
	HeldBalance float64 `gorm:"type:decimal(15,2);not null;default:0.00"`

	ExpiryMonth string `gorm:"size:2"`
	ExpiryYear  string `gorm:"size:4"`
	ExpiresAt   time.Time

	IssuedAt  time.Time `gorm:"not null;default:current_timestamp"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//
// =========================
// Transactions
// =========================
//

type Transaction struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	CardID uuid.UUID `gorm:"type:uuid;not null;index"`
	Card   Card      `gorm:"foreignKey:CardID"`

	TransactionReference string `gorm:"size:100;not null;uniqueIndex"`
	IdempotencyKey       *string `gorm:"size:100;index"`

	Amount   float64 `gorm:"type:decimal(15,2);not null"`
	Currency string  `gorm:"size:3;not null"`

	AuthorizedAmount float64 `gorm:"type:decimal(15,2)"`
	CapturedAmount   float64 `gorm:"type:decimal(15,2)"`

	Type      string `gorm:"size:50;not null"` // authorization, capture, funding, refund
	Direction string `gorm:"size:10;not null"` // debit | credit
	Status    string `gorm:"size:50;not null"`

	MerchantName    *string `gorm:"size:255"`
	MerchantMCC     *string `gorm:"size:4"`
	MerchantCountry *string `gorm:"size:2"`
	Source          *string `gorm:"size:30"` // card_network, bank_transfer

	DeclineReason *string `gorm:"type:text"`

	MetadataJSON datatypes.JSON `gorm:"type:jsonb"`


	TransactionTimestamp time.Time `gorm:"not null"`
	CreatedAt            time.Time `gorm:"autoCreateTime"`
}

//
// =========================
// Balance Ledger
// =========================
//

type BalanceLedger struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	CardID uuid.UUID `gorm:"type:uuid;not null;index"`
	Card   Card      `gorm:"foreignKey:CardID"`

	TransactionID uuid.UUID   `gorm:"type:uuid"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID"`

	EntryType    string  `gorm:"size:50"`
	Amount       float64 `gorm:"type:decimal(15,2);not null"`
	FeeCharged   float64 `gorm:"type:decimal(15,2);not null"`
	BalanceAfter float64 `gorm:"type:decimal(15,2);not null"`

	CreatedAt time.Time
}

//
// =========================
// Audit Logs
// =========================
//

type AuditLog struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID *uuid.UUID `gorm:"type:uuid;index"`
	User   *User      `gorm:"foreignKey:UserID"`

	Action     string `gorm:"size:100;not null"`
	EntityType string `gorm:"size:50;not null"`
	EntityID   *uuid.UUID `gorm:"type:uuid"`

	IPAddress  *string `gorm:"type:inet"`
	UserAgent  *string `gorm:"type:text"`
	RequestID  *string `gorm:"size:100"`
	Metadata   datatypes.JSON `gorm:"type:jsonb"`

	CreatedAt time.Time
}

//
// =========================
// Notifications
// =========================
//

type Notification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	Type    string `gorm:"size:50;not null"`
	Channel string `gorm:"size:20;not null"`

	Subject *string `gorm:"size:255"`
	Body    string  `gorm:"type:text;not null"`

	Status       string `gorm:"size:50;not null;default:pending"`
	RetryCount   int    `gorm:"not null;default:0"`
	ErrorMessage *string `gorm:"type:text"`

	SentAt    *time.Time
	CreatedAt time.Time
}

//
// =========================
// Refresh Tokens
// =========================
//

type RefreshToken struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	TokenHash string `gorm:"size:255;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`

	Revoked   bool       `gorm:"not null;default:false"`
	RevokedAt *time.Time

	CreatedAt time.Time
}
