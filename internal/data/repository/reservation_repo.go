package repository

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(reservation *entity.Reservation, tx ...*gorm.DB) error
	FindByID(id uint, tx ...*gorm.DB) (*entity.Reservation, error)
	FindAll(params ReservationQueryParams, tx ...*gorm.DB) ([]entity.Reservation, int64, error)
	FindByDate(date time.Time, tx ...*gorm.DB) ([]entity.Reservation, error)
	FindByCustomerID(customerID uint, tx ...*gorm.DB) ([]entity.Reservation, error)
	FindByTableID(tableID uint, tx ...*gorm.DB) ([]entity.Reservation, error)
	IsTableAvailable(tableID uint, date time.Time, reservationTime time.Time, tx ...*gorm.DB) (bool, error)
	GetReservationsByDateRange(startDate, endDate time.Time, tx ...*gorm.DB) ([]entity.Reservation, error)
	Update(reservation *entity.Reservation, tx ...*gorm.DB) error
	UpdateStatus(id uint, status entity.ReservationStatus, tx ...*gorm.DB) error
	Delete(id uint, tx ...*gorm.DB) error
}

type reservationRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewReservationRepository(db *gorm.DB, log *zap.Logger) ReservationRepository {
	return &reservationRepository{
		db:  db,
		log: log.With(zap.String("repository", "reservation")),
	}
}

type ReservationQueryParams struct {
	Date       *time.Time               // Filter by specific date
	StartDate  *time.Time               // Filter by date range start
	EndDate    *time.Time               // Filter by date range end
	Status     entity.ReservationStatus // Filter by status
	CustomerID *uint                    // Filter by customer
	TableID    *uint                    // Filter by table
	Offset     int                      // Pagination offset
	Limit      int                      // Pagination limit
}

// getDB returns either transaction db or regular db
func (rr *reservationRepository) getDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}
	return rr.db
}

func (rr *reservationRepository) Create(reservation *entity.Reservation, tx ...*gorm.DB) error {
	db := rr.getDB(tx...)

	rr.log.Debug("Creating new reservation",
		zap.Uint("customer_id", reservation.CustomerID),
		zap.Uint("table_id", reservation.TableID))

	err := db.Create(reservation).Error
	if err != nil {
		rr.log.Error("Failed to create reservation",
			zap.Uint("customer_id", reservation.CustomerID),
			zap.Uint("table_id", reservation.TableID),
			zap.Error(err))
		return err
	}

	rr.log.Info("Reservation created successfully",
		zap.Uint("id", reservation.ID),
		zap.Uint("customer_id", reservation.CustomerID))
	return nil
}

func (rr *reservationRepository) FindByID(id uint, tx ...*gorm.DB) (*entity.Reservation, error) {
	db := rr.getDB(tx...)

	rr.log.Debug("Finding reservation by ID", zap.Uint("id", id))

	var reservation entity.Reservation
	err := db.
		Preload("Customer").
		Preload("Table").
		First(&reservation, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rr.log.Warn("Reservation not found", zap.Uint("id", id))
		} else {
			rr.log.Error("Failed to find reservation by ID",
				zap.Uint("id", id),
				zap.Error(err))
		}
		return nil, err
	}

	rr.log.Debug("Reservation found", zap.Uint("id", id))
	return &reservation, nil
}

func (rr *reservationRepository) FindAll(params ReservationQueryParams, tx ...*gorm.DB) ([]entity.Reservation, int64, error) {
	db := rr.getDB(tx...)

	rr.log.Debug("Finding reservations",
		zap.Any("date", params.Date),
		zap.String("status", string(params.Status)),
		zap.Int("offset", params.Offset),
		zap.Int("limit", params.Limit))

	var reservations []entity.Reservation
	var total int64

	query := db.Model(&entity.Reservation{}).
		Preload("Customer").
		Preload("Table")

	// Apply date filter if provided
	if params.Date != nil {
		query = query.Where("reservation_date = ?", params.Date.Format("2006-01-02"))
	}

	// Apply date range filter if provided
	if params.StartDate != nil && params.EndDate != nil {
		query = query.Where("reservation_date BETWEEN ? AND ?",
			params.StartDate.Format("2006-01-02"),
			params.EndDate.Format("2006-01-02"))
	}

	// Apply status filter if provided
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Apply customer filter if provided
	if params.CustomerID != nil {
		query = query.Where("customer_id = ?", *params.CustomerID)
	}

	// Apply table filter if provided
	if params.TableID != nil {
		query = query.Where("table_id = ?", *params.TableID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		rr.log.Error("Failed to count reservations", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Offset(params.Offset).Limit(params.Limit)
	}

	// Execute query with ordering
	if err := query.
		Order("reservation_date DESC, reservation_time DESC").
		Find(&reservations).Error; err != nil {
		rr.log.Error("Failed to find reservations",
			zap.Error(err))
		return nil, 0, err
	}

	rr.log.Debug("Reservations found",
		zap.Int("count", len(reservations)),
		zap.Int64("total", total))
	return reservations, total, nil
}

func (rr *reservationRepository) FindByDate(date time.Time, tx ...*gorm.DB) ([]entity.Reservation, error) {
	db := rr.getDB(tx...)

	rr.log.Debug("Finding reservations by date", zap.Time("date", date))

	var reservations []entity.Reservation
	err := db.
		Preload("Customer").
		Preload("Table").
		Where("reservation_date = ?", date.Format("2006-01-02")).
		Order("reservation_time ASC").
		Find(&reservations).Error

	if err != nil {
		rr.log.Error("Failed to find reservations by date",
			zap.Time("date", date),
			zap.Error(err))
		return nil, err
	}

	rr.log.Debug("Reservations found by date",
		zap.Time("date", date),
		zap.Int("count", len(reservations)))
	return reservations, nil
}

func (rr *reservationRepository) FindByCustomerID(customerID uint, tx ...*gorm.DB) ([]entity.Reservation, error) {
	db := rr.getDB(tx...)

	rr.log.Debug("Finding reservations by customer", zap.Uint("customer_id", customerID))

	var reservations []entity.Reservation
	err := db.
		Preload("Customer").
		Preload("Table").
		Where("customer_id = ?", customerID).
		Order("reservation_date DESC").
		Find(&reservations).Error

	if err != nil {
		rr.log.Error("Failed to find reservations by customer",
			zap.Uint("customer_id", customerID),
			zap.Error(err))
		return nil, err
	}

	return reservations, nil
}

func (rr *reservationRepository) FindByTableID(tableID uint, tx ...*gorm.DB) ([]entity.Reservation, error) {
	db := rr.getDB(tx...)

	rr.log.Debug("Finding reservations by table", zap.Uint("table_id", tableID))

	var reservations []entity.Reservation
	err := db.
		Where("table_id = ?", tableID).
		Order("reservation_date DESC").
		Find(&reservations).Error

	if err != nil {
		rr.log.Error("Failed to find reservations by table",
			zap.Uint("table_id", tableID),
			zap.Error(err))
		return nil, err
	}

	return reservations, nil
}

func (rr *reservationRepository) IsTableAvailable(tableID uint, date time.Time, reservationTime time.Time, tx ...*gorm.DB) (bool, error) {
	db := rr.getDB(tx...)

	rr.log.Debug("Checking table availability",
		zap.Uint("table_id", tableID),
		zap.Time("date", date),
		zap.Time("time", reservationTime))

	var count int64

	dateStr := date.Format("2006-01-02")
	timeStr := reservationTime.Format("15:04:05")

	// Check for existing reservations that overlap with requested time
	// Consider reservations within ±2 hours as overlapping
	err := db.Model(&entity.Reservation{}).
		Where("table_id = ?", tableID).
		Where("reservation_date = ?", dateStr).
		Where("status IN ?", []entity.ReservationStatus{
			entity.ReservationStatusAwaiting,
			entity.ReservationStatusConfirmed,
		}).
		Where("ABS(EXTRACT(EPOCH FROM (reservation_time - ?::time))) / 3600 < 2", timeStr).
		Count(&count).Error

	if err != nil {
		rr.log.Error("Failed to check table availability",
			zap.Uint("table_id", tableID),
			zap.Time("date", date),
			zap.Error(err))
		return false, err
	}

	isAvailable := count == 0
	rr.log.Debug("Table availability check result",
		zap.Uint("table_id", tableID),
		zap.Bool("available", isAvailable),
		zap.Int64("conflicting_reservations", count))

	return isAvailable, nil
}

func (rr *reservationRepository) GetReservationsByDateRange(startDate, endDate time.Time, tx ...*gorm.DB) ([]entity.Reservation, error) {
	db := rr.getDB(tx...)

	rr.log.Debug("Getting reservations by date range",
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate))

	var reservations []entity.Reservation
	err := db.
		Preload("Customer").
		Preload("Table").
		Where("reservation_date BETWEEN ? AND ?",
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02")).
		Order("reservation_date ASC, reservation_time ASC").
		Find(&reservations).Error

	if err != nil {
		rr.log.Error("Failed to get reservations by date range",
			zap.Time("start_date", startDate),
			zap.Time("end_date", endDate),
			zap.Error(err))
		return nil, err
	}

	return reservations, err
}

func (rr *reservationRepository) Update(reservation *entity.Reservation, tx ...*gorm.DB) error {
	db := rr.getDB(tx...)

	rr.log.Debug("Updating reservation",
		zap.Uint("id", reservation.ID),
		zap.String("status", string(reservation.Status)))

	err := db.Save(reservation).Error
	if err != nil {
		rr.log.Error("Failed to update reservation",
			zap.Uint("id", reservation.ID),
			zap.Error(err))
		return err
	}

	rr.log.Info("Reservation updated", zap.Uint("id", reservation.ID))
	return nil
}

func (rr *reservationRepository) UpdateStatus(id uint, status entity.ReservationStatus, tx ...*gorm.DB) error {
	db := rr.getDB(tx...)

	rr.log.Debug("Updating reservation status",
		zap.Uint("id", id),
		zap.String("status", string(status)))

	err := db.Model(&entity.Reservation{}).
		Where("id = ?", id).
		Update("status", status).Error

	if err != nil {
		rr.log.Error("Failed to update reservation status",
			zap.Uint("id", id),
			zap.String("status", string(status)),
			zap.Error(err))
		return err
	}

	rr.log.Info("Reservation status updated",
		zap.Uint("id", id),
		zap.String("status", string(status)))
	return nil
}

func (rr *reservationRepository) Delete(id uint, tx ...*gorm.DB) error {
	db := rr.getDB(tx...)

	rr.log.Debug("Deleting reservation", zap.Uint("id", id))

	err := db.Delete(&entity.Reservation{}, id).Error
	if err != nil {
		rr.log.Error("Failed to delete reservation",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	rr.log.Info("Reservation deleted", zap.Uint("id", id))
	return nil
}
