package request

type CreateCustomerRequest struct {
	Title     string `json:"title" form:"title" validate:"omitempty,oneof=Mr Mrs Ms Dr Prof"`
	FirstName string `json:"first_name" form:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" form:"last_name"`
	Phone     string `json:"phone" form:"phone" validate:"required"`
	Email     string `json:"email" form:"email" validate:"omitempty,email"`
}

type UpdateCustomerRequest struct {
	Title     string `json:"title" form:"title" validate:"omitempty,oneof=Mr Mrs Ms Dr Prof"`
	FirstName string `json:"first_name" form:"first_name" validate:"omitempty,min=2"`
	LastName  string `json:"last_name" form:"last_name"`
	Phone     string `json:"phone" form:"phone"`
	Email     string `json:"email" form:"email" validate:"omitempty,email"`
}

type GetCustomersRequest struct {
	PaginationRequest
	Search string `json:"search" form:"search"`
}
