package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReportHandler struct {
	service usecase.ReportService
	log *zap.Logger
}

func NewReportHandler(service usecase.ReportService, log *zap.Logger) ReportHandler {
	return ReportHandler{
		service: service,
		log: log,
	}
}

func (h *ReportHandler) GetSales(c *gin.Context) {
	f := request.PeriodRequest{
		From: c.Query("from"),
		To: c.Query("to"),
	}

	result, err := h.service.GetSales(c, f)
	if err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "get revenue report failed", nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "get revenue report success", result)
}

func (h *ReportHandler) GetRevenue(c *gin.Context) {
	f := request.PeriodRequest{
		From: c.Query("from"),
		To: c.Query("to"),
	}

	result, err := h.service.GetRevenue(c, f)
	if err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "get revenue report failed", nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "get revenue report success", result)
}