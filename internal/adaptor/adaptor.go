package adaptor

import (
	"travel-api/internal/usecase"
	"travel-api/pkg/utils"

	"go.uber.org/zap"
)

type Handler struct{
	TourHandler TourHandler
}

func NewHandler(u *usecase.Usecase, log *zap.Logger, config utils.Configuration) Handler {
	return Handler{
		TourHandler: NewTourAdaptor(u, log, config),
	}
}