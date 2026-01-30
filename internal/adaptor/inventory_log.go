package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InventoryLogHandler struct {
	service usecase.InventoryLogService
	Logger *zap.Logger
	Config utils.Configuration
}

func NewInventoryLogHandler(service usecase.InventoryLogService, log *zap.Logger, config utils.Configuration) InventoryLogHandler {
	return InventoryLogHandler{
		service: service,
		Logger: log,
		Config: config,
	}
}

func (h *InventoryLogHandler) GetInventoryLogs(c *gin.Context) {
	var req request.InventoryLogsFilter
	result, err := h.service.GetInventoryLogs(c, req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadGateway, "get inventory logs failed", nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "get inventory logs success", result)
}

func (h *InventoryLogHandler) CreateInventoryLog(c *gin.Context) {
	var req request.CreateInventoryLogRequest
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

	res, err := h.service.CreateInventoryLog(c.Request.Context(), req)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "create inventory log failed", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "create inventory log success", res)
}