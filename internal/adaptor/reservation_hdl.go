package adaptor

import (
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/internal/usecase"
	"project-POS-APP-golang-integer/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReservationHandler struct {
	srv usecase.ReservationService
	log *zap.Logger
}

func NewReservationHandler(
	srv usecase.ReservationService,
	log *zap.Logger,
) *ReservationHandler {
	return &ReservationHandler{
		srv: srv,
		log: log.With(zap.String("handler", "reservation")),
	}
}

// CreateReservation creates a new reservation
func (rh *ReservationHandler) CreateReservation(c *gin.Context) {
	var req request.CreateReservationRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		rh.log.Warn("Failed to bind request body", zap.Error(err))
		utils.ResponseBadRequest(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		rh.log.Warn("Request validation failed",
			zap.Any("errors", validationErrors),
			zap.Error(err))

		utils.ResponseBadRequest(c, 400, "Validation failed", validationErrors)
		return
	}

	// Call service
	reservation, err := rh.srv.CreateReservation(req)
	if err != nil {
		rh.log.Error("Failed to create reservation", zap.Error(err))

		if utils.IsBusinessError(err) {
			utils.ResponseBadRequest(c, 400, err.Error(), nil)
		} else {
			utils.ResponseBadRequest(c, 500, "Failed to create reservation", nil)
		}
		return
	}

	rh.log.Info("Reservation created successfully",
		zap.Uint("reservation_id", reservation.ID))

	utils.ResponseSuccess(c, 201, "Reservation created successfully", reservation)
}

// GetReservations gets list of reservations with filters
func (rh *ReservationHandler) GetReservations(c *gin.Context) {
	var req request.GetReservationsRequest

	// Bind query
	if err := c.ShouldBindQuery(&req); err != nil {
		rh.log.Warn("Failed to bind query parameters", zap.Error(err))
		utils.ResponseBadRequest(c, 400, "Invalid query parameters", err.Error())
		return
	}

	// Call service
	result, err := rh.srv.GetReservations(req)
	if err != nil {
		rh.log.Error("Failed to get reservations", zap.Error(err))
		utils.ResponseBadRequest(c, 500, "Failed to get reservations", nil)
		return
	}

	rh.log.Debug("Reservations retrieved",
		zap.Int("count", len(result.Data)),
		zap.Int64("total", result.Pagination.Total))

	utils.ResponsePagination(c, 200, "Reservations retrieved successfully",
		result.Data, response.PaginationMeta{
			Page:       result.Pagination.Page,
			PerPage:    result.Pagination.PerPage,
			Total:      result.Pagination.Total,
			TotalPages: result.Pagination.TotalPages,
		})
}

// GetReservationByID gets a reservation by ID
func (rh *ReservationHandler) GetReservationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		rh.log.Warn("Invalid reservation ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseBadRequest(c, 400, "Invalid reservation ID", nil)
		return
	}

	// Call service
	reservation, err := rh.srv.GetReservationByID(uint(id))
	if err != nil {
		rh.log.Error("Failed to get reservation", zap.Uint("id", uint(id)), zap.Error(err))

		if err == utils.ErrReservationNotFound {
			utils.ResponseBadRequest(c, 404, "Reservation not found", nil)
		} else {
			utils.ResponseBadRequest(c, 500, "Failed to get reservation", nil)
		}
		return
	}

	utils.ResponseSuccess(c, 200, "Reservation retrieved successfully", reservation)
}

// UpdateReservationStatus updates reservation status
func (rh *ReservationHandler) UpdateReservationStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		rh.log.Warn("Invalid reservation ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseBadRequest(c, 400, "Invalid reservation ID", nil)
		return
	}

	var req request.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rh.log.Warn("Failed to bind request body", zap.Error(err))
		utils.ResponseBadRequest(c, 400, "Invalid request body", err.Error())
		return
	}

	// Validate
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		rh.log.Warn("Request validation failed", zap.Any("errors", validationErrors))
		utils.ResponseBadRequest(c, 400, "Validation failed", validationErrors)
		return
	}

	// Call service
	if err := rh.srv.UpdateReservationStatus(uint(id), req.Status); err != nil {
		rh.log.Error("Failed to update reservation status",
			zap.Uint("id", uint(id)),
			zap.String("status", req.Status),
			zap.Error(err))

		if err == utils.ErrReservationNotFound {
			utils.ResponseBadRequest(c, 404, "Reservation not found", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseBadRequest(c, 400, err.Error(), nil)
		} else {
			utils.ResponseBadRequest(c, 500, "Failed to update reservation status", nil)
		}
		return
	}

	rh.log.Info("Reservation status updated",
		zap.Uint("id", uint(id)),
		zap.String("status", req.Status))

	utils.ResponseSuccess(c, 200, "Reservation status updated successfully", nil)
}

// CancelReservation cancels a reservation
func (rh *ReservationHandler) CancelReservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		rh.log.Warn("Invalid reservation ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseBadRequest(c, 400, "Invalid reservation ID", nil)
		return
	}

	var req struct {
		Reason string `json:"reason" validate:"omitempty,max=500"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		rh.log.Warn("Failed to bind request body", zap.Error(err))
		utils.ResponseBadRequest(c, 400, "Invalid request body", err.Error())
		return
	}

	// Call service
	if err := rh.srv.CancelReservation(uint(id), req.Reason); err != nil {
		rh.log.Error("Failed to cancel reservation", zap.Uint("id", uint(id)), zap.Error(err))

		if err == utils.ErrReservationNotFound {
			utils.ResponseBadRequest(c, 404, "Reservation not found", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseBadRequest(c, 400, err.Error(), nil)
		} else {
			utils.ResponseBadRequest(c, 500, "Failed to cancel reservation", nil)
		}
		return
	}

	rh.log.Info("Reservation cancelled", zap.Uint("id", uint(id)))

	utils.ResponseSuccess(c, 200, "Reservation cancelled successfully", nil)
}

// CheckIn marks reservation as checked in
func (rh *ReservationHandler) CheckIn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		rh.log.Warn("Invalid reservation ID", zap.String("id", idStr), zap.Error(err))
		utils.ResponseBadRequest(c, 400, "Invalid reservation ID", nil)
		return
	}

	// Call service
	if err := rh.srv.CheckIn(uint(id)); err != nil {
		rh.log.Error("Failed to check in reservation", zap.Uint("id", uint(id)), zap.Error(err))

		if err == utils.ErrReservationNotFound {
			utils.ResponseBadRequest(c, 404, "Reservation not found", nil)
		} else if utils.IsBusinessError(err) {
			utils.ResponseBadRequest(c, 400, err.Error(), nil)
		} else {
			utils.ResponseBadRequest(c, 500, "Failed to check in reservation", nil)
		}
		return
	}

	rh.log.Info("Reservation checked in", zap.Uint("id", uint(id)))

	utils.ResponseSuccess(c, 200, "Reservation checked in successfully", nil)
}

// GetAvailableTables gets available tables for reservation
func (rh *ReservationHandler) GetAvailableTables(c *gin.Context) {
	date := c.Query("date")
	time := c.Query("time")
	paxStr := c.Query("pax")

	if date == "" || time == "" || paxStr == "" {
		utils.ResponseBadRequest(c, 400, "date, time, and pax parameters are required", nil)
		return
	}

	pax, err := strconv.Atoi(paxStr)
	if err != nil || pax <= 0 {
		utils.ResponseBadRequest(c, 400, "pax must be a positive integer", nil)
		return
	}

	// Call service
	tables, err := rh.srv.GetAvailableTables(date, time, pax)
	if err != nil {
		rh.log.Error("Failed to get available tables",
			zap.String("date", date),
			zap.String("time", time),
			zap.Int("pax", pax),
			zap.Error(err))

		if utils.IsBusinessError(err) {
			utils.ResponseBadRequest(c, 400, err.Error(), nil)
		} else {
			utils.ResponseBadRequest(c, 500, "Failed to get available tables", nil)
		}
		return
	}

	rh.log.Debug("Available tables found", zap.Int("count", len(tables)))

	utils.ResponseSuccess(c, 200, "Available tables retrieved successfully", tables)
}
