package response

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

type CustomerResponse struct {
	ID        uint                 `json:"id"`
	Title     entity.CustomerTitle `json:"title,omitempty"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name,omitempty"`
	Phone     string               `json:"phone"`
	Email     string               `json:"email,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type CustomerDetailResponse struct {
	CustomerResponse
	TotalReservations int `json:"total_reservations,omitempty"`
	TotalOrders       int `json:"total_orders,omitempty"`
}

// Converters
func CustomerToResponse(customer *entity.Customer) CustomerResponse {
	return CustomerResponse{
		ID:        customer.ID,
		Title:     customer.Title,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Phone:     customer.Phone,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

func CustomerToDetailResponse(customer *entity.Customer, totalReservations, totalOrders int) CustomerDetailResponse {
	return CustomerDetailResponse{
		CustomerResponse:  CustomerToResponse(customer),
		TotalReservations: totalReservations,
		TotalOrders:       totalOrders,
	}
}
