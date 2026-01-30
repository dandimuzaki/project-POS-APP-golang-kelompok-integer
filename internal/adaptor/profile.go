package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	service usecase.ProfileService
	Logger *zap.Logger
	Config utils.Configuration
}

func NewProfileHandler(service usecase.ProfileService, log *zap.Logger, config utils.Configuration) ProfileHandler {
	return ProfileHandler{
		service: service,
		Logger: log,
		Config: config,
	}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	result, err := h.service.GetProfile(c)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadGateway, "get profile failed", nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "get profile success", result)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req request.ProfileRequest
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

	err = h.service.UpdateProfile(c, &req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "update profile failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "update profile success", nil)
}