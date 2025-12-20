package models


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