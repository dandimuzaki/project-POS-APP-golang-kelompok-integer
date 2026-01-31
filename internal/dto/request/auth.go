package request

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"email"`
}

type ResetPassword struct {
	NewPassword string `json:"new_password"`
	ResetToken  string `json:"reset_token"`
}

type ValidateOTP struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}