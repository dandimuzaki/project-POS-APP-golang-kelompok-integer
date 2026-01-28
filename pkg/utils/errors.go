package utils

import "errors"

// =============== ERROR RESERVATION ===============
var (
	ErrTableUnavailable        = errors.New("table is not available at selected time")
	ErrInvalidReservationTime  = errors.New("reservation time must be at least 1 hour from now")
	ErrCustomerNotFound        = errors.New("customer not found")
	ErrTableNotFound           = errors.New("table not found")
	ErrReservationNotFound     = errors.New("reservation not found")
	ErrInsufficientCapacity    = errors.New("no tables available with sufficient capacity")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// Helper untuk check business error
func IsBusinessError(err error) bool {
	switch err {
	case ErrTableUnavailable,
		ErrInvalidReservationTime,
		ErrCustomerNotFound,
		ErrTableNotFound,
		ErrReservationNotFound,
		ErrInsufficientCapacity,
		ErrInvalidStatusTransition:
		return true
	default:
		return false
	}
}
