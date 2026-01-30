package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service usecase.AuthService
	Logger  *zap.Logger
	Config  utils.Configuration
}

func NewAuthHandler(service usecase.AuthService, log *zap.Logger, config utils.Configuration) AuthHandler {
	return AuthHandler{
		service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), messages)
		return
	}

	res, err := h.service.Login(c, req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusUnauthorized, "login failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "login success", res)
}

func (h *AuthHandler) RequestResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), messages)
		return
	}

	result, err := h.service.RequestResetPassword(c, req.Email)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "request reset password failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "login success", result)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), messages)
		return
	}

	err = h.service.ResetPassword(c, req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "reset password failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "reset password success", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	
}