package utils

import "errors"

var (
	// =============== ERROR AUTH ===============
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidOTP = errors.New("invalid OTP")
	ErrInvalidToken = errors.New("invalid token")

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

	// =============== ERROR CATEGORY ===============
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategoryExists      = errors.New("category name already exists")
	ErrCategoryHasProducts = errors.New("cannot delete category with associated products")
	ErrInvalidCategoryName = errors.New("category name is invalid")
	ErrCategoryInactive    = errors.New("category is inactive")
	ErrNoChangesProvided   = errors.New("no changes provided") // 🔥 Optional

	// =============== ERROR PRODUCT ===============
	ErrProductNotFound   = errors.New("product not found")
	ErrProductExists     = errors.New("product name already exists in category")
	ErrProductHasOrders  = errors.New("cannot delete product with associated orders")
	ErrProductInactive   = errors.New("product is inactive")
	ErrProductOutOfStock = errors.New("product is out of stock")
	ErrInsufficientStock = errors.New("insufficient stock")
)

// Helper untuk check business error
func IsBusinessError(err error) bool {
	businessErrors := []error{
		// Reservation errors
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

		// Customer errors
		ErrCustomerNotFound,
		ErrCustomerAlreadyExists,

		// Category errors
		ErrCategoryNotFound,
		ErrCategoryExists,
		ErrCategoryHasProducts,
		ErrInvalidCategoryName,
		ErrCategoryInactive,
		ErrNoChangesProvided, // 🔥
		// Product errors
		ErrProductNotFound,
		ErrProductExists,
		ErrProductHasOrders,
		ErrProductInactive,
		ErrProductOutOfStock,
		ErrInsufficientStock,
	}

	for _, businessErr := range businessErrors {
		if errors.Is(err, businessErr) {
			return true
		}
	}
	return false
}
