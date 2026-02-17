package models

import (
	"time"

	"github.com/google/uuid"
)


type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPwdReq struct {
	Email string `json:"email"`
}

type ForgotPwdOtp struct{
	Otp string  `json:"otp"`
	Email string `json:"email"`
}

type ResetPasswordReq struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordReq struct{
	Userid uuid.UUID
	OldPassword string `json:"old_password"`
	NewPassword string 	`json:"new_password"`
}

type MFALoginReq struct  {
	Email    string `json:"email"`
	TOTPCode string `json:"totp_code"`
}


type Otp struct{
	Otp string `json:"otp"`
}

type VerifyMFA struct{
	TotpCode string `json:"totp_code"`
}

type KycProfile struct{
	Userid uuid.UUID
	DOB string `json:"dob"`
	ImageStr string `json:"image_string"`
}

type KycDoc struct{
	Userid uuid.UUID
	DocStr string `json:"doc_string"`
}


type CreateCardReq struct{
	Userid uuid.UUID
	CardType string `json:"card_type"`
	Currency string `json:"currency"`
	SpendingLimit float64 `json:"spending_limit"`
}

type CreateCardResp struct{
	CardType string `json:"card_type"`
	MaskedPAN string `json:"masked_pan"`
	Currency string `json:"currency"`
	SpendingLimit float64 `json:"spending_limit"`
	Balance float64 `json:"balance"`
	Cvv string `json:"cvv"`
	Status string `json:"status"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear string `json:"expiry_year"`

}

type GetAllCardsResp struct {
	Cardid uuid.UUID `json:"card_id"`
	CardType string `json:"card_type"`
	MaskedPAN string `json:"masked_pan"`
	Lastfour string `json:"last_four"`
	Cvv string `json:"cvv"`
	Currency string `json:"currency"`
	Status string `json:"status"`
	
}

type GetCardReq struct{
	UserId uuid.UUID
	CardId string
}

type GetCardResp struct{
	Cardid uuid.UUID `json:"card_id"`
	CardType string `json:"card_type"`
	PAN string `json:"card_number"`
	Cvv string `json:"cvv"`
	Lastfour string `json:"last_four"`
	Currency string `json:"currency"`
	Status string `json:"status"`
	SpendingLimit float64 `json:"spending_limit"`
	CurrentBalance float64 `json:"current_balance"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear string `json:"expiry_year"`
}

type StatusReq struct{
	Status string `json:"status"`
}

type CardEmailPayload struct{
	Email string
	FirstName string
	LastFour string
	Type string
}

type TopUpCardReq struct{
	Userid uuid.UUID
	Cardid string `json:"card_id"`
	Amount float64 `json:"amount"`
}

type TopUpCardResp struct{
	AccountName string  `json:"account_name"`
	AccountNumber string `json:"account_number"`
	Bank string `json:"bank"`
	Reference string `json:"reference"`
	Note string `json:"note"`
}

type merchant struct{
	Name string `json:"name"`
	MCC string `json:"mcc"`
	Country string `json:"country"`
}
type WebhookReq struct { 
	OriginalTransactionID string `json:"original_transaction_id"` //used for card refund events
	TransactionID string `json:"transaction_id"`
	CardReference string `json:"card_reference"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency"`
	Type string `json:"type"`
	Direction string `json:"direction"`
	Status string `json:"status"`
	Merchant merchant `json:"merchant"`
	Network string `json:"network"`
	Timestamp time.Time `json:"timestamp"`
	IdempotencyKey string `json:"idempotency_key"`
}

type GetCardTransactionsReq struct {
	Userid uuid.UUID
	Cardid string `json:"card_id"`
}

type GetCardTransactionsResp struct{
	Cardid uuid.UUID `json:"card_id"`
	Transaction_Reference string `json:"transaction_reference"`
	Amount float64 `json:"amount"`
	AuthorizedAmount float64 `json:"authorized _amount"`
	CapturedAmount float64 `json:"captured_amount"`
	Currency string `json:"currency"`
	MerchantName *string `json:"merchant_name"`
	Direction string `json:"direction"`
	Type string  `json:"type"`
	Source *string `json:"source"`
	DeclineReason *string `json:"decline_reason"`
	CreatedAt time.Time `json:"created_at"`
}