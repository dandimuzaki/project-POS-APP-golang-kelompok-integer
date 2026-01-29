package request

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
}

type ResetPassword struct {
	Email       string `json:"email" validate:"email"`
	NewPassword string `json:"new_password"`
	OTP         string `json:"otp"`
}