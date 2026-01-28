package usecase

import (
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
	CreateReservation(req request.CreateReservationRequest) (*response.ReservationResponse, error)
	GetReservations(req request.GetReservationsRequest) (*response.PaginatedResponse[response.ReservationResponse], error)
	GetReservationByID(id uint) (*response.ReservationResponse, error)
	UpdateReservationStatus(id uint, status string) error
	CancelReservation(id uint, reason string) error
	CheckIn(reservationID uint) error
	GetAvailableTables(dateStr, timeStr string, paxNumber int) ([]response.TableResponse, error)
}

type reservationService struct {
	db   *gorm.DB
	repo *repository.Repository
	log  *zap.Logger
}

func NewReservationService(
	db *gorm.DB,
	repo *repository.Repository,
	log *zap.Logger,
) ReservationService {
	return &reservationService{
		db:   db,
		repo: repo,
		log:  log.With(zap.String("service", "reservation")),
	}
}

func (rs *reservationService) CreateReservation(req request.CreateReservationRequest) (*response.ReservationResponse, error) {
	rs.log.Info("Creating new reservation",
		zap.String("customer_name", req.Customer.FirstName),
		zap.Int("pax_number", req.Reservation.PaxNumber))

	// 1. Validate input
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		rs.log.Warn("Validation failed for reservation request",
			zap.Any("errors", validationErrors))
		return nil, err
	}

	// 2. Parse dates
	reservationDate, err := utils.ParseReservationDate(req.Reservation.ReservationDate)
	if err != nil {
		rs.log.Error("Failed to parse reservation date",
			zap.String("date", req.Reservation.ReservationDate),
			zap.Error(err))
		return nil, err
	}

	reservationTime, err := utils.ParseReservationTime(req.Reservation.ReservationTime)
	if err != nil {
		rs.log.Error("Failed to parse reservation time",
			zap.String("time", req.Reservation.ReservationTime),
			zap.Error(err))
		return nil, err
	}

	// 3. Validate reservation time (must be at least 1 hour from now)
	if !utils.IsValidReservationTime(reservationDate, reservationTime) {
		rs.log.Warn("Invalid reservation time",
			zap.Time("date", reservationDate),
			zap.Time("time", reservationTime))
		return nil, utils.ErrInvalidReservationTime
	}

	// 4. START TRANSACTION
	tx := rs.db.Begin()
	if tx.Error != nil {
		rs.log.Error("Failed to begin transaction", zap.Error(tx.Error))
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			rs.log.Error("Transaction panic, rolled back", zap.Any("panic", r))
		}
	}()

	// 5. Find or create customer (WITHIN TRANSACTION)
	customer, err := rs.repo.CustomerRepo.FindByPhone(req.Customer.Phone, tx)
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		rs.log.Error("Failed to find customer", zap.Error(err))
		return nil, err
	}

	if customer == nil {
		rs.log.Debug("Customer not found, creating new one",
			zap.String("phone", req.Customer.Phone))

		customer = &entity.Customer{
			Title:     entity.CustomerTitle(req.Customer.Title),
			FirstName: req.Customer.FirstName,
			LastName:  req.Customer.LastName,
			Phone:     req.Customer.Phone,
			Email:     req.Customer.Email,
		}

		if err := rs.repo.CustomerRepo.Create(customer, tx); err != nil {
			tx.Rollback()
			rs.log.Error("Failed to create customer",
				zap.String("phone", req.Customer.Phone),
				zap.Error(err))
			return nil, err
		}

		rs.log.Info("New customer created",
			zap.Uint("customer_id", customer.ID),
			zap.String("name", customer.FirstName))
	} else {
		rs.log.Debug("Existing customer found",
			zap.Uint("customer_id", customer.ID),
			zap.String("name", customer.FirstName))
	}

	// 6. Find available table (WITHIN TRANSACTION)
	var table *entity.Table
	if req.Reservation.TableID > 0 {
		// Customer specified a table
		table, err = rs.repo.TableRepo.FindByID(req.Reservation.TableID, tx)
		if err != nil {
			tx.Rollback()
			rs.log.Error("Specified table not found",
				zap.Uint("table_id", req.Reservation.TableID),
				zap.Error(err))
			return nil, utils.ErrTableNotFound
		}

		// Check if table is available
		isAvailable, err := rs.repo.ReservationRepo.IsTableAvailable(
			table.ID, reservationDate, reservationTime, tx,
		)
		if err != nil || !isAvailable {
			tx.Rollback()
			rs.log.Warn("Specified table is not available",
				zap.Uint("table_id", table.ID),
				zap.Time("date", reservationDate),
				zap.Time("time", reservationTime))
			return nil, utils.ErrTableUnavailable
		}

		// Check capacity
		if table.Capacity < req.Reservation.PaxNumber {
			tx.Rollback()
			rs.log.Warn("Table capacity insufficient",
				zap.Uint("table_id", table.ID),
				zap.Int("capacity", table.Capacity),
				zap.Int("pax_number", req.Reservation.PaxNumber))
			return nil, utils.ErrInsufficientCapacity
		}
	} else {
		// Auto-select table based on capacity
		tables, err := rs.repo.TableRepo.FindByCapacity(req.Reservation.PaxNumber, tx)
		if err != nil || len(tables) == 0 {
			tx.Rollback()
			rs.log.Error("No tables available with sufficient capacity",
				zap.Int("pax_number", req.Reservation.PaxNumber),
				zap.Error(err))
			return nil, utils.ErrInsufficientCapacity
		}

		// Find first available table
		for _, t := range tables {
			isAvailable, _ := rs.repo.ReservationRepo.IsTableAvailable(
				t.ID, reservationDate, reservationTime, tx,
			)
			if isAvailable {
				table = &t
				break
			}
		}

		if table == nil {
			tx.Rollback()
			rs.log.Warn("No tables available at selected time",
				zap.Time("date", reservationDate),
				zap.Time("time", reservationTime),
				zap.Int("pax_number", req.Reservation.PaxNumber))
			return nil, utils.ErrTableUnavailable
		}
	}

	// 7. Create reservation (WITHIN TRANSACTION)
	reservation := &entity.Reservation{
		CustomerID:      customer.ID,
		TableID:         table.ID,
		PaxNumber:       req.Reservation.PaxNumber,
		ReservationDate: reservationDate,
		ReservationTime: reservationTime,
		Status:          entity.ReservationStatusAwaiting,
		Notes:           req.Reservation.Notes,
	}

	if err := rs.repo.ReservationRepo.Create(reservation, tx); err != nil {
		tx.Rollback()
		rs.log.Error("Failed to create reservation",
			zap.Uint("customer_id", customer.ID),
			zap.Uint("table_id", table.ID),
			zap.Error(err))
		return nil, err
	}

	// 8. Update table status (WITHIN TRANSACTION)
	if err := rs.repo.TableRepo.UpdateStatus(table.ID, entity.TableStatusReserved, tx); err != nil {
		tx.Rollback()
		rs.log.Error("Failed to update table status",
			zap.Uint("table_id", table.ID),
			zap.Error(err))
		return nil, err
	}

	// 9. COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		rs.log.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	rs.log.Info("Reservation transaction committed successfully",
		zap.Uint("reservation_id", reservation.ID),
		zap.String("customer", customer.FirstName),
		zap.String("table", table.TableNumber))

	// 10. Load complete data for response (outside transaction)
	reservationWithDetails, err := rs.repo.ReservationRepo.FindByID(reservation.ID)
	if err != nil {
		rs.log.Error("Failed to load reservation details",
			zap.Uint("reservation_id", reservation.ID),
			zap.Error(err))
		// Don't fail, just return basic reservation
		return &response.ReservationResponse{
			ID:              reservation.ID,
			Customer:        response.CustomerToResponse(customer),
			Table:           response.TableToResponse(table),
			PaxNumber:       reservation.PaxNumber,
			ReservationDate: utils.FormatReservationDate(reservationDate),
			ReservationTime: utils.FormatReservationTime(reservationTime),
			Status:          reservation.Status,
			Notes:           reservation.Notes,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}, nil
	}

	resp := response.ReservationToResponse(reservationWithDetails)
	return &resp, nil
}

func (rs *reservationService) GetReservations(req request.GetReservationsRequest) (*response.PaginatedResponse[response.ReservationResponse], error) {
	rs.log.Debug("Getting reservations",
		zap.String("date", req.Date),
		zap.String("status", req.Status),
		zap.Int("page", req.GetPage()),
		zap.Int("per_page", req.GetPerPage()))

	// Convert request to repo params
	params := repository.ReservationQueryParams{
		Offset: req.GetOffset(),
		Limit:  req.GetPerPage(),
		Status: entity.ReservationStatus(req.Status),
	}

	// Parse date if provided
	if req.Date != "" {
		date, err := utils.ParseReservationDate(req.Date)
		if err != nil {
			rs.log.Warn("Invalid date format",
				zap.String("date", req.Date),
				zap.Error(err))
			return nil, err
		}
		params.Date = &date
	}

	// Apply customer filter
	if req.CustomerID > 0 {
		params.CustomerID = &req.CustomerID
	}

	// Apply table filter
	if req.TableID > 0 {
		params.TableID = &req.TableID
	}

	// Get reservations from repository (WITHOUT TRANSACTION for read)
	reservations, total, err := rs.repo.ReservationRepo.FindAll(params)
	if err != nil {
		rs.log.Error("Failed to get reservations",
			zap.Error(err))
		return nil, err
	}

	// Convert to DTOs
	reservationDTOs := make([]response.ReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		reservationDTOs = append(reservationDTOs, response.ReservationToResponse(&r))
	}

	rs.log.Debug("Reservations retrieved",
		zap.Int("count", len(reservationDTOs)),
		zap.Int64("total", total))

	return response.NewPaginatedResponse(
		reservationDTOs,
		req.GetPage(),
		req.GetPerPage(),
		total,
	), nil
}

func (rs *reservationService) GetReservationByID(id uint) (*response.ReservationResponse, error) {
	rs.log.Debug("Getting reservation by ID", zap.Uint("id", id))

	reservation, err := rs.repo.ReservationRepo.FindByID(id)
	if err != nil {
		rs.log.Error("Failed to get reservation",
			zap.Uint("id", id),
			zap.Error(err))
		return nil, utils.ErrReservationNotFound
	}

	resp := response.ReservationToResponse(reservation)
	return &resp, nil
}

func (rs *reservationService) UpdateReservationStatus(id uint, status string) error {
	rs.log.Info("Updating reservation status",
		zap.Uint("id", id),
		zap.String("status", status))

	// Validate status
	reservationStatus := entity.ReservationStatus(status)
	if !reservationStatus.IsValid() {
		rs.log.Warn("Invalid reservation status",
			zap.String("status", status))
		return utils.ErrInvalidStatusTransition
	}

	// START TRANSACTION
	tx := rs.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get reservation within transaction
	reservation, err := rs.repo.ReservationRepo.FindByID(id, tx)
	if err != nil {
		tx.Rollback()
		return utils.ErrReservationNotFound
	}

	// Validate status transition (business rule)
	if !rs.isValidStatusTransition(reservation.Status, reservationStatus) {
		tx.Rollback()
		rs.log.Warn("Invalid status transition",
			zap.String("from", string(reservation.Status)),
			zap.String("to", string(reservationStatus)))
		return utils.ErrInvalidStatusTransition
	}

	// Update reservation status within transaction
	if err := rs.repo.ReservationRepo.UpdateStatus(id, reservationStatus, tx); err != nil {
		tx.Rollback()
		rs.log.Error("Failed to update reservation status",
			zap.Uint("id", id),
			zap.String("status", status),
			zap.Error(err))
		return err
	}

	// If cancelled or completed, free the table
	if reservationStatus == entity.ReservationStatusCancelled ||
		reservationStatus == entity.ReservationStatusCompleted {
		if err := rs.repo.TableRepo.UpdateStatus(reservation.TableID, entity.TableStatusAvailable, tx); err != nil {
			tx.Rollback()
			rs.log.Error("Failed to free table",
				zap.Uint("table_id", reservation.TableID),
				zap.Error(err))
			return err
		}
	}

	// COMMIT
	if err := tx.Commit().Error; err != nil {
		return err
	}

	rs.log.Info("Reservation status updated successfully",
		zap.Uint("id", id),
		zap.String("from", string(reservation.Status)),
		zap.String("to", string(reservationStatus)))

	return nil
}

// Helper method for status transition validation
func (rs *reservationService) isValidStatusTransition(from, to entity.ReservationStatus) bool {
	// Define valid transitions
	validTransitions := map[entity.ReservationStatus][]entity.ReservationStatus{
		entity.ReservationStatusAwaiting: {
			entity.ReservationStatusConfirmed,
			entity.ReservationStatusCancelled,
		},
		entity.ReservationStatusConfirmed: {
			entity.ReservationStatusCompleted,
			entity.ReservationStatusCancelled,
		},
		entity.ReservationStatusCompleted: {}, // No transitions from completed
		entity.ReservationStatusCancelled: {}, // No transitions from cancelled
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

func (rs *reservationService) CancelReservation(id uint, reason string) error {
	rs.log.Info("Cancelling reservation",
		zap.Uint("id", id),
		zap.String("reason", reason))

	// START TRANSACTION
	tx := rs.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get reservation within transaction
	reservation, err := rs.repo.ReservationRepo.FindByID(id, tx)
	if err != nil {
		tx.Rollback()
		return utils.ErrReservationNotFound
	}

	// Validate status (can only cancel awaiting or confirmed)
	if reservation.Status != entity.ReservationStatusAwaiting &&
		reservation.Status != entity.ReservationStatusConfirmed {
		tx.Rollback()
		rs.log.Warn("Cannot cancel reservation in current status",
			zap.String("current_status", string(reservation.Status)))
		return utils.ErrInvalidStatusTransition
	}

	// Update reservation status to cancelled
	reservation.Status = entity.ReservationStatusCancelled
	if reason != "" {
		reservation.Notes += "\nCancellation reason: " + reason
	}

	if err := rs.repo.ReservationRepo.Update(reservation, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Free the table
	if err := rs.repo.TableRepo.UpdateStatus(reservation.TableID, entity.TableStatusAvailable, tx); err != nil {
		tx.Rollback()
		return err
	}

	// COMMIT
	if err := tx.Commit().Error; err != nil {
		return err
	}

	rs.log.Info("Reservation cancelled successfully",
		zap.Uint("id", id),
		zap.String("reason", reason))

	return nil
}

func (rs *reservationService) CheckIn(reservationID uint) error {
	rs.log.Info("Checking in reservation", zap.Uint("reservation_id", reservationID))

	// Get reservation
	reservation, err := rs.repo.ReservationRepo.FindByID(reservationID)
	if err != nil {
		return utils.ErrReservationNotFound
	}

	// Validate status
	if reservation.Status != entity.ReservationStatusConfirmed {
		rs.log.Warn("Cannot check in reservation in current status",
			zap.String("status", string(reservation.Status)))
		return utils.ErrInvalidStatusTransition
	}

	// Update check out time (when customer arrives)
	now := time.Now()
	reservation.CheckOutAt = &now

	if err := rs.repo.ReservationRepo.Update(reservation); err != nil {
		return err
	}

	rs.log.Info("Reservation checked in successfully",
		zap.Uint("reservation_id", reservationID))

	return nil
}

func (rs *reservationService) GetAvailableTables(dateStr, timeStr string, paxNumber int) ([]response.TableResponse, error) {
	rs.log.Debug("Getting available tables",
		zap.String("date", dateStr),
		zap.String("time", timeStr),
		zap.Int("pax_number", paxNumber))

	// Parse dates
	reservationDate, err := utils.ParseReservationDate(dateStr)
	if err != nil {
		return nil, err
	}

	reservationTime, err := utils.ParseReservationTime(timeStr)
	if err != nil {
		return nil, err
	}

	// Get tables with sufficient capacity
	tables, err := rs.repo.TableRepo.FindByCapacity(paxNumber)
	if err != nil {
		return nil, err
	}

	// Filter tables that are available at requested time
	var availableTables []entity.Table
	for _, table := range tables {
		isAvailable, err := rs.repo.ReservationRepo.IsTableAvailable(
			table.ID, reservationDate, reservationTime,
		)
		if err == nil && isAvailable {
			availableTables = append(availableTables, table)
		}
	}

	// Convert to DTOs
	tableDTOs := make([]response.TableResponse, 0, len(availableTables))
	for _, table := range availableTables {
		tableDTOs = append(tableDTOs, response.TableToResponse(&table))
	}

	rs.log.Debug("Available tables found",
		zap.Int("count", len(tableDTOs)))

	return tableDTOs, nil
}
