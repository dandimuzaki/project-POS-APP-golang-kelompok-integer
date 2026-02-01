package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentMethodHandler struct {
	service usecase.PaymentMethodService
	log *zap.Logger
}

func NewPaymentMethodHandler(service usecase.PaymentMethodService, log *zap.Logger) PaymentMethodHandler {
	return PaymentMethodHandler{
		service: service,
		log: log,
	}
}

func (h *PaymentMethodHandler) GetAll(c *gin.Context) {
	result, err := h.service.GetAll(c)
	if err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "get payment methods failed", nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "get payment methods success", result)
}
