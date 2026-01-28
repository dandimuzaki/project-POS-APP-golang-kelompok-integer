package response

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

type ReservationResponse struct {
	ID              uint                     `json:"id"`
	Customer        CustomerResponse         `json:"customer"`
	Table           TableResponse            `json:"table"`
	PaxNumber       int                      `json:"pax_number"`
	ReservationDate string                   `json:"reservation_date"`
	ReservationTime string                   `json:"reservation_time"`
	DepositFee      float64                  `json:"deposit_fee"`
	Status          entity.ReservationStatus `json:"status"`
	Notes           string                   `json:"notes,omitempty"`
	CheckOutAt      *time.Time               `json:"check_out_at,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

type ReservationDetailResponse struct {
	ReservationResponse
}

// Converters
func ReservationToResponse(reservation *entity.Reservation) ReservationResponse {
	// Format dates
	reservationDate := ""
	if !reservation.ReservationDate.IsZero() {
		reservationDate = reservation.ReservationDate.Format("2006-01-02")
	}

	reservationTime := ""
	if !reservation.ReservationTime.IsZero() {
		reservationTime = reservation.ReservationTime.Format("15:04")
	}

	return ReservationResponse{
		ID:              reservation.ID,
		Customer:        CustomerToResponse(&reservation.Customer),
		Table:           TableToResponse(&reservation.Table),
		PaxNumber:       reservation.PaxNumber,
		ReservationDate: reservationDate,
		ReservationTime: reservationTime,
		DepositFee:      reservation.DepositFee,
		Status:          reservation.Status,
		Notes:           reservation.Notes,
		CheckOutAt:      reservation.CheckOutAt,
		CreatedAt:       reservation.CreatedAt,
		UpdatedAt:       reservation.UpdatedAt,
	}
}

func ReservationToDetailResponse(reservation *entity.Reservation) ReservationDetailResponse {
	return ReservationDetailResponse{
		ReservationResponse: ReservationToResponse(reservation),
	}
}
