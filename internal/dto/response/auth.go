package response

import (
	"time"

	"github.com/google/uuid"
)

type AuthResponse struct {
	Token uuid.UUID `json:"token"`
}

type OTPResponse struct {
	OTPCode   string `json:"otp_code"`
	ExpiresAt time.Time `json:"expires_at"`
}