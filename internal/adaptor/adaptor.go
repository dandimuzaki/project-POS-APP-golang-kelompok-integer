package adaptor

import (
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
)

type Handler struct {
	UserHandler         UserHandler
	AuthHandler         AuthHandler
	ProfileHandler      ProfileHandler
	ReservationHandler  ReservationHandler
	InventoryLogHandler InventoryLogHandler
	CategoryHandler     CategoryHandler
	ProductHandler      ProductHandler
}

func NewHandler(u *usecase.Usecase, log *zap.Logger, config utils.Configuration) Handler {
	return Handler{
		UserHandler:         NewUserHandler(u.UserService, log, config),
		AuthHandler:         NewAuthHandler(u.AuthService, log, config),
		ProfileHandler:      NewProfileHandler(u.ProfileService, log, config),
		ReservationHandler:  NewReservationHandler(u.ReservationService, log, config),
		InventoryLogHandler: NewInventoryLogHandler(u.InventoryLogService, log, config),
		CategoryHandler:     *NewCategoryHandler(u.CategoryService, log),
		ProductHandler:      *NewProductHandler(u.ProductService, log),
	}
}
