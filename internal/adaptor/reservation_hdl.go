package adaptor

import (
	"net/http"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReservationHandler struct {
	service usecase.ReservationService
	logger  *zap.Logger
	config  utils.Configuration
}

func NewReservationHandler(service usecase.ReservationService, log *zap.Logger, config utils.Configuration) ReservationHandler {
	return ReservationHandler{
		service: service,
		logger:  log.With(zap.String("handler", "reservation")),
		config:  config,
	}
}

// CreateReservation creates a new reservation
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req request.CreateReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path))
		utils.ResponseFailed(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate request menggunakan utils dari kamu
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		h.logger.Warn("Validation failed",
			zap.Any("errors", validationErrors))
		utils.ResponseFailed(c, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	reservation, err := h.service.CreateReservation(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create reservation",
			zap.Error(err),
			zap.String("customer_phone", req.Customer.Phone))

		if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, http.StatusInternalServerError, "Failed to create reservation", nil)
		}
		return
	}

	h.logger.Info("Reservation created successfully",
		zap.Uint("reservation_id", reservation.ID),
		zap.String("customer_name", reservation.Customer.FirstName))

	utils.ResponseSuccess(c, http.StatusCreated, "Reservation created successfully", reservation)
}

// GetReservations gets list of reservations
func (h *ReservationHandler) GetReservations(c *gin.Context) {
	var req request.GetReservationsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Warn("Invalid query parameters",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path))
		utils.ResponseFailed(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	reservations, pagination, err := h.service.GetReservations(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get reservations",
			zap.Error(err),
			zap.Any("filters", req))
		utils.ResponseFailed(c, http.StatusInternalServerError, "Failed to get reservations", nil)
		return
	}

	h.logger.Debug("Reservations retrieved",
		zap.Int("count", len(reservations)),
		zap.Int64("total", pagination.Total))

	utils.ResponsePagination(c, http.StatusOK, "Reservations retrieved successfully",
		reservations, pagination)
}

// GetReservationByID gets a reservation by ID
func (h *ReservationHandler) GetReservationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid reservation ID",
			zap.String("id", idStr),
			zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "Invalid reservation ID", nil)
		return
	}

	reservation, err := h.service.GetReservationByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get reservation",
			zap.Uint("id", uint(id)),
			zap.Error(err))

		if err == utils.ErrReservationNotFound {
			utils.ResponseFailed(c, http.StatusNotFound, "Reservation not found", nil)
		} else {
			utils.ResponseFailed(c, http.StatusInternalServerError, "Failed to get reservation", nil)
		}
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Reservation retrieved successfully", reservation)
}

// UpdateReservationStatus updates reservation status
func (h *ReservationHandler) UpdateReservationStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid reservation ID",
			zap.String("id", idStr),
			zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "Invalid reservation ID", nil)
		return
	}

	var req request.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path))
		utils.ResponseFailed(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		h.logger.Warn("Validation failed",
			zap.Any("errors", validationErrors))
		utils.ResponseFailed(c, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	if err := h.service.UpdateReservationStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		h.logger.Error("Failed to update reservation status",
			zap.Uint("id", uint(id)),
			zap.String("status", req.Status),
			zap.Error(err))

		if err == utils.ErrReservationNotFound {
			utils.ResponseFailed(c, http.StatusNotFound, "Reservation not found", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, http.StatusInternalServerError, "Failed to update reservation status", nil)
		}
		return
	}

	h.logger.Info("Reservation status updated",
		zap.Uint("id", uint(id)),
		zap.String("status", req.Status))

	utils.ResponseSuccess(c, http.StatusOK, "Reservation status updated successfully", nil)
}

// CancelReservation cancels a reservation
func (h *ReservationHandler) CancelReservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid reservation ID",
			zap.String("id", idStr),
			zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "Invalid reservation ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason" validate:"omitempty,max=500"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body",
			zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.service.CancelReservation(c.Request.Context(), uint(id), req.Reason); err != nil {
		h.logger.Error("Failed to cancel reservation",
			zap.Uint("id", uint(id)),
			zap.String("reason", req.Reason),
			zap.Error(err))

		if err == utils.ErrReservationNotFound {
			utils.ResponseFailed(c, http.StatusNotFound, "Reservation not found", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, http.StatusInternalServerError, "Failed to cancel reservation", nil)
		}
		return
	}

	h.logger.Info("Reservation cancelled",
		zap.Uint("id", uint(id)),
		zap.String("reason", req.Reason))

	utils.ResponseSuccess(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

// CheckIn checks in a reservation
func (h *ReservationHandler) CheckIn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid reservation ID",
			zap.String("id", idStr),
			zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "Invalid reservation ID", nil)
		return
	}

	if err := h.service.CheckIn(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to check in reservation",
			zap.Uint("id", uint(id)),
			zap.Error(err))

		if err == utils.ErrReservationNotFound {
			utils.ResponseFailed(c, http.StatusNotFound, "Reservation not found", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, http.StatusInternalServerError, "Failed to check in reservation", nil)
		}
		return
	}

	h.logger.Info("Reservation checked in", zap.Uint("id", uint(id)))
	utils.ResponseSuccess(c, http.StatusOK, "Reservation checked in successfully", nil)
}

// GetAvailableTables gets available tables
func (h *ReservationHandler) GetAvailableTables(c *gin.Context) {
	date := c.Query("date")
	time := c.Query("time")
	paxStr := c.Query("pax")

	if date == "" || time == "" || paxStr == "" {
		h.logger.Warn("Missing required parameters",
			zap.String("date", date),
			zap.String("time", time),
			zap.String("pax", paxStr))
		utils.ResponseFailed(c, http.StatusBadRequest, "date, time, and pax parameters are required", nil)
		return
	}

	pax, err := strconv.Atoi(paxStr)
	if err != nil || pax <= 0 {
		h.logger.Warn("Invalid pax parameter",
			zap.String("pax", paxStr),
			zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "pax must be a positive integer", nil)
		return
	}

	tables, err := h.service.GetAvailableTables(c.Request.Context(), date, time, pax)
	if err != nil {
		h.logger.Error("Failed to get available tables",
			zap.String("date", date),
			zap.String("time", time),
			zap.Int("pax", pax),
			zap.Error(err))

		if utils.IsBusinessError(err) {
			utils.ResponseFailed(c, http.StatusBadRequest, err.Error(), nil)
		} else {
			utils.ResponseFailed(c, http.StatusInternalServerError, "Failed to get available tables", nil)
		}
		return
	}

	h.logger.Debug("Available tables retrieved",
		zap.Int("count", len(tables)))

	utils.ResponseSuccess(c, http.StatusOK, "Available tables retrieved successfully", tables)
}
