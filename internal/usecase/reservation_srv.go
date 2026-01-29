package usecase

import (
	"context"
	"math"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/pkg/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReservationService interface {
	CreateReservation(ctx context.Context, req request.CreateReservationRequest) (*response.ReservationResponse, error)
	GetReservations(ctx context.Context, req request.GetReservationsRequest) ([]response.ReservationResponse, response.PaginationMeta, error)
	GetReservationByID(ctx context.Context, id uint) (*response.ReservationResponse, error)
	UpdateReservationStatus(ctx context.Context, id uint, status string) error
	CancelReservation(ctx context.Context, id uint, reason string) error
	CheckIn(ctx context.Context, id uint) error
	GetAvailableTables(ctx context.Context, dateStr, timeStr string, paxNumber int) ([]response.TableResponse, error)
}

type reservationService struct {
	tx   TxManager
	repo *repository.Repository
	log  *zap.Logger
}

func NewReservationService(
	tx TxManager,
	repo *repository.Repository,
	log *zap.Logger,
) ReservationService {
	return &reservationService{
		tx:   tx,
		repo: repo,
		log:  log.With(zap.String("service", "reservation")),
	}
}

func (s *reservationService) CreateReservation(ctx context.Context, req request.CreateReservationRequest) (*response.ReservationResponse, error) {
	s.log.Info("Creating new reservation",
		zap.String("customer_name", req.Customer.FirstName),
		zap.Int("pax_number", req.Reservation.PaxNumber))

	// 1. Validate request
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		s.log.Warn("Validation failed",
			zap.Any("errors", validationErrors))
		return nil, utils.ErrValidationFailed
	}

	// 2. Parse dates
	reservationDate, err := time.Parse("2006-01-02", req.Reservation.ReservationDate)
	if err != nil {
		s.log.Error("Invalid reservation date format",
			zap.String("date", req.Reservation.ReservationDate),
			zap.Error(err))
		return nil, utils.ErrInvalidDateFormat
	}

	reservationTime, err := time.Parse("15:04", req.Reservation.ReservationTime)
	if err != nil {
		s.log.Error("Invalid reservation time format",
			zap.String("time", req.Reservation.ReservationTime),
			zap.Error(err))
		return nil, utils.ErrInvalidTimeFormat
	}

	// 3. Validate reservation time
	if !s.isValidReservationTime(reservationDate, reservationTime) {
		s.log.Warn("Invalid reservation time",
			zap.Time("date", reservationDate),
			zap.Time("time", reservationTime))
		return nil, utils.ErrInvalidReservationTime
	}

	var reservation *entity.Reservation
	var customer *entity.Customer
	var table *entity.Table

	// 4. Execute in transaction using TxManager
	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// 5. Find or create customer
		customer, err = s.repo.CustomerRepo.FindByPhone(ctx, req.Customer.Phone)
		if err != nil && err != gorm.ErrRecordNotFound {
			s.log.Error("Failed to find customer",
				zap.String("phone", req.Customer.Phone),
				zap.Error(err))
			return err
		}

		if customer == nil {
			s.log.Debug("Creating new customer",
				zap.String("phone", req.Customer.Phone))

			customer = &entity.Customer{
				Title:     entity.CustomerTitle(req.Customer.Title),
				FirstName: req.Customer.FirstName,
				LastName:  req.Customer.LastName,
				Phone:     req.Customer.Phone,
				Email:     req.Customer.Email,
			}

			customer, err = s.repo.CustomerRepo.Create(ctx, customer)
			if err != nil {
				s.log.Error("Failed to create customer",
					zap.String("phone", req.Customer.Phone),
					zap.Error(err))
				return err
			}

			s.log.Info("New customer created",
				zap.Uint("customer_id", customer.ID),
				zap.String("name", customer.FirstName))
		} else {
			s.log.Debug("Existing customer found",
				zap.Uint("customer_id", customer.ID),
				zap.String("name", customer.FirstName))
		}

		// 6. Find or select table
		if req.Reservation.TableID > 0 {
			// Customer specified a table
			table, err = s.repo.TableRepo.FindByID(ctx, req.Reservation.TableID)
			if err != nil {
				s.log.Error("Specified table not found",
					zap.Uint("table_id", req.Reservation.TableID),
					zap.Error(err))
				return utils.ErrTableNotFound
			}

			// Check availability
			isAvailable, err := s.repo.ReservationRepo.IsTableAvailable(
				ctx, table.ID, reservationDate, reservationTime)
			if err != nil || !isAvailable {
				s.log.Warn("Table not available",
					zap.Uint("table_id", table.ID),
					zap.Time("date", reservationDate),
					zap.Time("time", reservationTime))
				return utils.ErrTableUnavailable
			}

			// Check capacity
			if table.Capacity < req.Reservation.PaxNumber {
				s.log.Warn("Table capacity insufficient",
					zap.Uint("table_id", table.ID),
					zap.Int("capacity", table.Capacity),
					zap.Int("pax", req.Reservation.PaxNumber))
				return utils.ErrInsufficientCapacity
			}
		} else {
			// Auto-select table
			tables, err := s.repo.TableRepo.FindByCapacity(ctx, req.Reservation.PaxNumber)
			if err != nil || len(tables) == 0 {
				s.log.Error("No tables available with sufficient capacity",
					zap.Int("pax", req.Reservation.PaxNumber),
					zap.Error(err))
				return utils.ErrInsufficientCapacity
			}

			// Find first available table
			for _, t := range tables {
				isAvailable, _ := s.repo.ReservationRepo.IsTableAvailable(
					ctx, t.ID, reservationDate, reservationTime)
				if isAvailable {
					table = &t
					break
				}
			}

			if table == nil {
				s.log.Warn("No tables available at selected time",
					zap.Time("date", reservationDate),
					zap.Time("time", reservationTime),
					zap.Int("pax", req.Reservation.PaxNumber))
				return utils.ErrTableUnavailable
			}
		}

		// 7. Create reservation
		reservation = &entity.Reservation{
			CustomerID:      customer.ID,
			TableID:         table.ID,
			PaxNumber:       req.Reservation.PaxNumber,
			ReservationDate: reservationDate,
			ReservationTime: reservationTime,
			Status:          entity.ReservationStatusAwaiting,
			Notes:           req.Reservation.Notes,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		reservation, err = s.repo.ReservationRepo.Create(ctx, reservation)
		if err != nil {
			s.log.Error("Failed to create reservation",
				zap.Uint("customer_id", customer.ID),
				zap.Uint("table_id", table.ID),
				zap.Error(err))
			return err
		}

		// 8. Update table status
		if err := s.repo.TableRepo.UpdateStatus(ctx, table.ID, entity.TableStatusReserved); err != nil {
			s.log.Error("Failed to update table status",
				zap.Uint("table_id", table.ID),
				zap.Error(err))
			return err
		}

		s.log.Info("Reservation created within transaction",
			zap.Uint("reservation_id", reservation.ID),
			zap.String("customer", customer.FirstName),
			zap.String("table", table.TableNumber))

		return nil
	})

	if err != nil {
		s.log.Error("Reservation transaction failed", zap.Error(err))
		return nil, err
	}

	// 9. Load complete data for response
	reservationWithDetails, err := s.repo.ReservationRepo.FindByID(ctx, reservation.ID)
	if err != nil {
		s.log.Error("Failed to load reservation details",
			zap.Uint("reservation_id", reservation.ID),
			zap.Error(err))
		// Return basic info if loading fails
		return &response.ReservationResponse{
			ID:              reservation.ID,
			Customer:        response.CustomerToResponse(customer),
			Table:           response.TableToResponse(table),
			PaxNumber:       reservation.PaxNumber,
			ReservationDate: reservationDate.Format("2006-01-02"),
			ReservationTime: reservationTime.Format("15:04"),
			Status:          reservation.Status,
			Notes:           reservation.Notes,
			CreatedAt:       reservation.CreatedAt,
			UpdatedAt:       reservation.UpdatedAt,
		}, nil
	}

	resp := response.ReservationToResponse(reservationWithDetails)
	s.log.Info("Reservation created successfully",
		zap.Uint("reservation_id", reservation.ID),
		zap.String("customer", customer.FirstName),
		zap.String("table", table.TableNumber))

	return &resp, nil
}

func (s *reservationService) GetReservations(ctx context.Context, req request.GetReservationsRequest) ([]response.ReservationResponse, response.PaginationMeta, error) {
	s.log.Debug("Getting reservations",
		zap.String("date", req.Date),
		zap.String("status", req.Status),
		zap.Int("page", req.GetPage()),
		zap.Int("per_page", req.GetPerPage()))

	reservations, total, err := s.repo.ReservationRepo.FindAll(ctx, req)
	if err != nil {
		s.log.Error("Failed to get reservations", zap.Error(err))
		return nil, response.PaginationMeta{}, err
	}

	// Convert to DTOs
	reservationDTOs := make([]response.ReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		reservationDTOs = append(reservationDTOs, response.ReservationToResponse(&r))
	}

	// Calculate pagination
	totalPages := 0
	if req.GetPerPage() > 0 && total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(req.GetPerPage())))
	}

	pagination := response.PaginationMeta{
		Page:       req.GetPage(),
		PerPage:    req.GetPerPage(),
		Total:      total,
		TotalPages: totalPages,
	}

	s.log.Debug("Reservations retrieved",
		zap.Int("count", len(reservationDTOs)),
		zap.Int64("total", total))

	return reservationDTOs, pagination, nil
}

func (s *reservationService) GetReservationByID(ctx context.Context, id uint) (*response.ReservationResponse, error) {
	s.log.Debug("Getting reservation by ID", zap.Uint("id", id))

	reservation, err := s.repo.ReservationRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.log.Warn("Reservation not found", zap.Uint("id", id))
			return nil, utils.ErrReservationNotFound
		}
		s.log.Error("Failed to get reservation",
			zap.Uint("id", id),
			zap.Error(err))
		return nil, err
	}

	resp := response.ReservationToResponse(reservation)
	return &resp, nil
}

func (s *reservationService) UpdateReservationStatus(ctx context.Context, id uint, status string) error {
	s.log.Info("Updating reservation status",
		zap.Uint("id", id),
		zap.String("status", status))

	// Validate status
	reservationStatus := entity.ReservationStatus(status)
	if !reservationStatus.IsValid() {
		s.log.Warn("Invalid reservation status", zap.String("status", status))
		return utils.ErrInvalidStatus
	}

	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// Get current reservation
		reservation, err := s.repo.ReservationRepo.FindByID(ctx, id)
		if err != nil {
			s.log.Error("Reservation not found", zap.Uint("id", id), zap.Error(err))
			return utils.ErrReservationNotFound
		}

		// Validate status transition
		if !s.isValidStatusTransition(reservation.Status, reservationStatus) {
			s.log.Warn("Invalid status transition",
				zap.String("from", string(reservation.Status)),
				zap.String("to", string(reservationStatus)))
			return utils.ErrInvalidStatusTransition
		}

		// Update status
		if err := s.repo.ReservationRepo.UpdateStatus(ctx, id, reservationStatus); err != nil {
			s.log.Error("Failed to update reservation status",
				zap.Uint("id", id),
				zap.Error(err))
			return err
		}

		// If cancelled or completed, free the table
		if reservationStatus == entity.ReservationStatusCancelled ||
			reservationStatus == entity.ReservationStatusCompleted {
			if err := s.repo.TableRepo.UpdateStatus(ctx, reservation.TableID, entity.TableStatusAvailable); err != nil {
				s.log.Error("Failed to free table",
					zap.Uint("table_id", reservation.TableID),
					zap.Error(err))
				return err
			}
		}

		s.log.Info("Reservation status updated successfully",
			zap.Uint("id", id),
			zap.String("from", string(reservation.Status)),
			zap.String("to", string(reservationStatus)))

		return nil
	})
}

func (s *reservationService) CancelReservation(ctx context.Context, id uint, reason string) error {
	s.log.Info("Cancelling reservation",
		zap.Uint("id", id),
		zap.String("reason", reason))

	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// Get reservation
		reservation, err := s.repo.ReservationRepo.FindByID(ctx, id)
		if err != nil {
			return utils.ErrReservationNotFound
		}

		// Validate status
		if reservation.Status != entity.ReservationStatusAwaiting &&
			reservation.Status != entity.ReservationStatusConfirmed {
			s.log.Warn("Cannot cancel reservation in current status",
				zap.String("current_status", string(reservation.Status)))
			return utils.ErrInvalidStatusTransition
		}

		// Update to cancelled
		reservation.Status = entity.ReservationStatusCancelled
		if reason != "" {
			reservation.Notes += "\nCancellation reason: " + reason
		}

		if err := s.repo.ReservationRepo.Update(ctx, reservation); err != nil {
			s.log.Error("Failed to cancel reservation",
				zap.Uint("id", id),
				zap.Error(err))
			return err
		}

		// Free the table
		if err := s.repo.TableRepo.UpdateStatus(ctx, reservation.TableID, entity.TableStatusAvailable); err != nil {
			s.log.Error("Failed to free table",
				zap.Uint("table_id", reservation.TableID),
				zap.Error(err))
			return err
		}

		s.log.Info("Reservation cancelled successfully",
			zap.Uint("id", id),
			zap.String("reason", reason))

		return nil
	})
}

func (s *reservationService) CheckIn(ctx context.Context, id uint) error {
	s.log.Info("Checking in reservation", zap.Uint("id", id))

	return s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// Get reservation
		reservation, err := s.repo.ReservationRepo.FindByID(ctx, id)
		if err != nil {
			return utils.ErrReservationNotFound
		}

		// Validate status
		if reservation.Status != entity.ReservationStatusConfirmed {
			s.log.Warn("Cannot check in reservation in current status",
				zap.String("status", string(reservation.Status)))
			return utils.ErrInvalidStatusTransition
		}

		// Update check in time
		now := time.Now()
		reservation.CheckOutAt = &now

		// Update table status to occupied
		if err := s.repo.TableRepo.UpdateStatus(ctx, reservation.TableID, entity.TableStatusOccupied); err != nil {
			s.log.Error("Failed to update table status",
				zap.Uint("table_id", reservation.TableID),
				zap.Error(err))
			return err
		}

		if err := s.repo.ReservationRepo.Update(ctx, reservation); err != nil {
			s.log.Error("Failed to check in reservation",
				zap.Uint("id", id),
				zap.Error(err))
			return err
		}

		s.log.Info("Reservation checked in successfully", zap.Uint("id", id))
		return nil
	})
}

func (s *reservationService) GetAvailableTables(ctx context.Context, dateStr, timeStr string, paxNumber int) ([]response.TableResponse, error) {
	s.log.Debug("Getting available tables",
		zap.String("date", dateStr),
		zap.String("time", timeStr),
		zap.Int("pax_number", paxNumber))

	// Parse dates
	reservationDate, err := utils.ParseReservationDate(dateStr)
	if err != nil {
		s.log.Error("Invalid date format", zap.String("date", dateStr), zap.Error(err))
		return nil, utils.ErrInvalidDateFormat
	}

	reservationTime, err := utils.ParseReservationTime(timeStr)
	if err != nil {
		s.log.Error("Invalid time format", zap.String("time", timeStr), zap.Error(err))
		return nil, utils.ErrInvalidTimeFormat
	}

	// Get tables with sufficient capacity
	tables, err := s.repo.TableRepo.FindByCapacity(ctx, paxNumber)
	if err != nil {
		s.log.Error("Failed to get tables by capacity",
			zap.Int("pax", paxNumber),
			zap.Error(err))
		return nil, err
	}

	// Filter available tables
	var availableTables []entity.Table
	for _, table := range tables {
		isAvailable, err := s.repo.ReservationRepo.IsTableAvailable(
			ctx, table.ID, reservationDate, reservationTime)
		if err == nil && isAvailable {
			availableTables = append(availableTables, table)
		}
	}

	// Convert to DTOs
	tableDTOs := make([]response.TableResponse, 0, len(availableTables))
	for _, table := range availableTables {
		tableDTOs = append(tableDTOs, response.TableToResponse(&table))
	}

	s.log.Debug("Available tables found", zap.Int("count", len(tableDTOs)))
	return tableDTOs, nil
}

// Helper methods
func (s *reservationService) isValidReservationTime(date time.Time, reservationTime time.Time) bool {
	now := time.Now()
	reservationDateTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		reservationTime.Hour(), reservationTime.Minute(), 0, 0, time.Local)

	// Must be at least 1 hour from now
	minAllowedTime := now.Add(1 * time.Hour)
	return reservationDateTime.After(minAllowedTime)
}

func (s *reservationService) isValidStatusTransition(from, to entity.ReservationStatus) bool {
	validTransitions := map[entity.ReservationStatus][]entity.ReservationStatus{
		entity.ReservationStatusAwaiting: {
			entity.ReservationStatusConfirmed,
			entity.ReservationStatusCancelled,
		},
		entity.ReservationStatusConfirmed: {
			entity.ReservationStatusCompleted,
			entity.ReservationStatusCancelled,
		},
		entity.ReservationStatusCompleted: {},
		entity.ReservationStatusCancelled: {},
	}

	allowedTransitions, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == to {
			return true
		}
	}

	return false
}
