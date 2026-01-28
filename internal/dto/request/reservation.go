package request

type CreateReservationRequest struct {
	Customer struct {
		Title     string `json:"title" form:"title" validate:"omitempty,oneof=Mr Mrs Ms Dr Prof"`
		FirstName string `json:"first_name" form:"first_name" validate:"required,min=2"`
		LastName  string `json:"last_name" form:"last_name"`
		Phone     string `json:"phone" form:"phone" validate:"required"`
		Email     string `json:"email" form:"email" validate:"omitempty,email"`
	} `json:"customer"`

	Reservation struct {
		PaxNumber       int    `json:"pax_number" form:"pax_number" validate:"required,min=1,max=20"`
		ReservationDate string `json:"reservation_date" form:"reservation_date" validate:"required"`
		ReservationTime string `json:"reservation_time" form:"reservation_time" validate:"required"`
		TableID         uint   `json:"table_id" form:"table_id" validate:"omitempty"`
		Notes           string `json:"notes" form:"notes"`
	} `json:"reservation"`
}

type UpdateReservationRequest struct {
	Status  string `json:"status" form:"status" validate:"omitempty,oneof=awaiting confirmed cancelled completed"`
	TableID uint   `json:"table_id" form:"table_id" validate:"omitempty"`
	Notes   string `json:"notes" form:"notes"`
}

type GetReservationsRequest struct {
	PaginationRequest
	Date       string `json:"date" form:"date"`
	Status     string `json:"status" form:"status"`
	CustomerID uint   `json:"customer_id" form:"customer_id"`
	TableID    uint   `json:"table_id" form:"table_id"`
}
