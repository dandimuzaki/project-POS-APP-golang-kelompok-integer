package email

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/dto/response"
)

func AccountCreated(data response.CreateUserResponse) string {
	return fmt.Sprintf(`
	<h2>Your Account Has Been Created</h2>

	<p>
	An administrator has created an account for you.
	You can use the credentials below to log in:
	</p>

	<table style="
		border-collapse: collapse;
		width: %s;
		margin: 16px 0;
	">
		<tr>
			<td style="padding: 8px; font-weight: bold;">Email</td>
			<td style="padding: 8px;">%v</td>
		</tr>
		<tr>
			<td style="padding: 8px; font-weight: bold;">Temporary Password</td>
			<td style="padding: 8px;">%v</td>
		</tr>
	</table>

	<p>
	For security reasons, please <strong>change your password immediately</strong>
	after logging in.
	</p>

	<p>
	If you did not expect this account, please contact support.
	</p>

	<p style="color: #888; font-size: 12px;">
	⚠️ Do not share your login credentials with anyone.
	</p>
	`, "100%", data.Email, data.Password)
}