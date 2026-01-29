package adaptor

import (
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
)

type Handler struct {
	UserHandler        UserHandler
	AuthHandler        AuthHandler
	ReservationHandler ReservationHandler
}

func NewHandler(u *usecase.Usecase, log *zap.Logger, config utils.Configuration) Handler {
	return Handler{
		UserHandler:        NewUserHandler(u, log, config),
		AuthHandler:        NewAuthHandler(u, log, config),
		ReservationHandler: NewReservationHandler(u.ReservationService, log, config),
	}
}
