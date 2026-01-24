package middleware

import (
	"project-POS-APP-golang-integer/internal/usecase"

	"go.uber.org/zap"
)

type MiddlewareCustom struct {
	Usecase *usecase.Usecase
	Log     *zap.Logger
}

func NewMiddlewareCustom(usecase *usecase.Usecase, log *zap.Logger) MiddlewareCustom {
	return MiddlewareCustom{
		Usecase: usecase,
		Log:     log,
	}
}
