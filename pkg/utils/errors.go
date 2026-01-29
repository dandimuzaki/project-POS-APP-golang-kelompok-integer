package utils

import "errors"

var (
	// =============== ERROR RESERVATION ===============
	ErrReservationNotFound     = errors.New("reservation not found")
	ErrTableNotFound           = errors.New("table not found")
	ErrTableUnavailable        = errors.New("table is not available at the requested time")
	ErrInsufficientCapacity    = errors.New("table capacity is insufficient")
	ErrInvalidReservationTime  = errors.New("reservation time must be at least 1 hour from now")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidStatus           = errors.New("invalid reservation status")
	ErrValidationFailed        = errors.New("validation failed")
	ErrInvalidDateFormat       = errors.New("invalid date format")
	ErrInvalidTimeFormat       = errors.New("invalid time format")

	// Customer errors
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrCustomerAlreadyExists = errors.New("customer already exists")
)

func IsBusinessError(err error) bool {
	businessErrors := []error{
		ErrReservationNotFound,
		ErrTableNotFound,
		ErrTableUnavailable,
		ErrInsufficientCapacity,
		ErrInvalidReservationTime,
		ErrInvalidStatus,
		ErrInvalidStatusTransition,
		ErrValidationFailed,
		ErrInvalidDateFormat,
		ErrInvalidTimeFormat,
		ErrCustomerNotFound,
		ErrCustomerAlreadyExists,
	}

	for _, businessErr := range businessErrors {
		if errors.Is(err, businessErr) {
			return true
		}
	}
	return false
}
