package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"time"
)

func GenerateTransactionID(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)
	return encoded[:length], nil
}

// GenerateOTP generates a cryptographically secure numeric OTP of a given length.
func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		b[i] = digits[num.Int64()]
	}
	return string(b), nil
}

// =============== RESERVATION DATETIME HELPER ===============
// Parse reservation date string (YYYY-MM-DD) to time.Time
func ParseReservationDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// Parse reservation time string (HH:MM) to time.Time
func ParseReservationTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

// Check if reservation time is valid (not in the past)
func IsValidReservationTime(date time.Time, reservationTime time.Time) bool {
	now := time.Now()
	reservationDateTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		reservationTime.Hour(), reservationTime.Minute(), 0, 0, time.Local,
	)

	// Reservation must be at least 1 hour from now
	minAllowedTime := now.Add(1 * time.Hour)
	return reservationDateTime.After(minAllowedTime)
}

// Format time for display
func FormatReservationTime(t time.Time) string {
	return t.Format("15:04")
}

func FormatReservationDate(t time.Time) string {
	return t.Format("2006-01-02")
}
