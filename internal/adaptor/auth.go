package adaptor

import (
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service *usecase.Usecase
	Logger  *zap.Logger
	Config  utils.Configuration
}

func NewAuthHandler(service *usecase.Usecase, log *zap.Logger, config utils.Configuration) AuthHandler {
	return AuthHandler{
		service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	
}