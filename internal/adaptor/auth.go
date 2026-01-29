package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/dto"
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
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	res, err := h.service.AuthService.Login(c.Request.Context(), req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusUnauthorized, "login failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "login success", res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	
}