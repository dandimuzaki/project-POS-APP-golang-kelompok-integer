package email

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/dto/response"
)

func SendOTP(data response.OTPResponse) string {
	return fmt.Sprintf(`
	<h2>Password Reset Request</h2>

	<p>
	We received a request to reset your password.
	Please use the OTP code below to continue:
	</p>

	<div style="
		font-size: 26px;
		font-weight: bold;
		letter-spacing: 6px;
		margin: 20px 0;
		text-align: center;
	">
		%v
	</div>

	<p>
	This OTP will expire in <strong>10 minutes</strong>.
	</p>

	<p>
	If you did not request a password reset, please ignore this email.
	Your account remains secure.
	</p>

	<p style="color: #888; font-size: 12px;">
	⚠️ Do not share this OTP with anyone.
	</p>`, data.OTPCode)
}