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

// =============== ERROR CATEGORY ===============
var (
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategoryExists      = errors.New("category name already exists")
	ErrCategoryHasProducts = errors.New("cannot delete category with associated products")
	ErrInvalidCategoryName = errors.New("category name is invalid")
	ErrCategoryInactive    = errors.New("category is inactive")
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
		ErrInvalidStatusTransition,
		ErrCategoryNotFound,
		ErrCategoryExists,
		ErrCategoryHasProducts,
		ErrInvalidCategoryName,
		ErrCategoryInactive:
		return true
	default:
		return false
	}
}
