package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/infra"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error)
	FindByID(ctx context.Context, id uint) (*entity.Reservation, error)
	FindAll(ctx context.Context, params request.GetReservationsRequest) ([]entity.Reservation, int64, error)
	FindByCustomerID(ctx context.Context, customerID uint) ([]entity.Reservation, error)
	FindByDate(ctx context.Context, date time.Time) ([]entity.Reservation, error)
	IsTableAvailable(ctx context.Context, tableID uint, date time.Time, reservationTime time.Time) (bool, error)
	Update(ctx context.Context, reservation *entity.Reservation) error
	UpdateStatus(ctx context.Context, id uint, status entity.ReservationStatus) error
	Delete(ctx context.Context, id uint) error
}

type reservationRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewReservationRepo(db *gorm.DB, log *zap.Logger) ReservationRepository {
	return &reservationRepository{
		db:     db,
		logger: log.With(zap.String("repository", "reservation")),
	}
}

func (r *reservationRepository) Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Creating reservation",
		zap.Uint("customer_id", reservation.CustomerID),
		zap.Uint("table_id", reservation.TableID),
		zap.Time("date", reservation.ReservationDate),
		zap.String("status", string(reservation.Status)))

	err := db.Create(reservation).Error
	if err != nil {
		r.logger.Error("Failed to create reservation",
			zap.Error(err),
			zap.Uint("customer_id", reservation.CustomerID))
		return nil, err
	}

	r.logger.Info("Reservation created successfully",
		zap.Uint("id", reservation.ID),
		zap.Uint("customer_id", reservation.CustomerID))

	return reservation, nil
}

func (r *reservationRepository) FindByID(ctx context.Context, id uint) (*entity.Reservation, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding reservation by ID", zap.Uint("id", id))

	var reservation entity.Reservation
	err := db.
		Preload("Customer").
		Preload("Table").
		First(&reservation, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn("Reservation not found", zap.Uint("id", id))
		} else {
			r.logger.Error("Failed to find reservation",
				zap.Uint("id", id),
				zap.Error(err))
		}
		return nil, err
	}

	r.logger.Debug("Reservation found", zap.Uint("id", id))
	return &reservation, nil
}

func (r *reservationRepository) FindAll(ctx context.Context, params request.GetReservationsRequest) ([]entity.Reservation, int64, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding reservations",
		zap.String("date", params.Date),
		zap.String("status", params.Status),
		zap.Int("page", params.GetPage()),
		zap.Int("per_page", params.GetPerPage()))

	var reservations []entity.Reservation
	var total int64

	query := db.Model(&entity.Reservation{}).
		Preload("Customer").
		Preload("Table")

	// Apply filters
	if params.Date != "" {
		date, err := time.Parse("2006-01-02", params.Date)
		if err == nil {
			query = query.Where("DATE(reservation_date) = ?", date.Format("2006-01-02"))
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.CustomerID > 0 {
		query = query.Where("customer_id = ?", params.CustomerID)
	}

	if params.TableID > 0 {
		query = query.Where("table_id = ?", params.TableID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count reservations", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	offset := params.GetOffset()
	limit := params.GetPerPage()

	err := query.
		Offset(offset).
		Limit(limit).
		Order("reservation_date DESC, reservation_time DESC").
		Find(&reservations).Error

	if err != nil {
		r.logger.Error("Failed to find reservations",
			zap.Error(err),
			zap.Int("offset", offset),
			zap.Int("limit", limit))
		return nil, 0, err
	}

	r.logger.Debug("Reservations retrieved",
		zap.Int("count", len(reservations)),
		zap.Int64("total", total))

	return reservations, total, nil
}

func (r *reservationRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]entity.Reservation, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding reservations by customer", zap.Uint("customer_id", customerID))

	var reservations []entity.Reservation
	err := db.
		Where("customer_id = ?", customerID).
		Order("reservation_date DESC").
		Find(&reservations).Error

	if err != nil {
		r.logger.Error("Failed to find reservations by customer",
			zap.Uint("customer_id", customerID),
			zap.Error(err))
		return nil, err
	}

	return reservations, nil
}

func (r *reservationRepository) FindByDate(ctx context.Context, date time.Time) ([]entity.Reservation, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding reservations by date", zap.Time("date", date))

	var reservations []entity.Reservation
	err := db.
		Where("DATE(reservation_date) = ?", date.Format("2006-01-02")).
		Order("reservation_time ASC").
		Find(&reservations).Error

	if err != nil {
		r.logger.Error("Failed to find reservations by date",
			zap.Time("date", date),
			zap.Error(err))
		return nil, err
	}

	return reservations, nil
}

func (r *reservationRepository) IsTableAvailable(ctx context.Context, tableID uint, date time.Time, reservationTime time.Time) (bool, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Checking table availability",
		zap.Uint("table_id", tableID),
		zap.Time("date", date),
		zap.Time("time", reservationTime))

	var count int64
	dateStr := date.Format("2006-01-02")
	timeStr := reservationTime.Format("15:04:05")

	err := db.Model(&entity.Reservation{}).
		Where("table_id = ?", tableID).
		Where("DATE(reservation_date) = ?", dateStr).
		Where("status IN ?", []entity.ReservationStatus{
			entity.ReservationStatusAwaiting,
			entity.ReservationStatusConfirmed,
		}).
		Where("ABS(EXTRACT(EPOCH FROM (reservation_time - ?::time))) / 3600 < 2", timeStr).
		Count(&count).Error

	if err != nil {
		r.logger.Error("Failed to check table availability",
			zap.Uint("table_id", tableID),
			zap.Error(err))
		return false, err
	}

	isAvailable := count == 0
	r.logger.Debug("Table availability result",
		zap.Uint("table_id", tableID),
		zap.Bool("available", isAvailable),
		zap.Int64("conflicting", count))

	return isAvailable, nil
}

func (r *reservationRepository) Update(ctx context.Context, reservation *entity.Reservation) error {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Updating reservation",
		zap.Uint("id", reservation.ID),
		zap.String("status", string(reservation.Status)))

	err := db.Save(reservation).Error
	if err != nil {
		r.logger.Error("Failed to update reservation",
			zap.Uint("id", reservation.ID),
			zap.Error(err))
		return err
	}

	r.logger.Info("Reservation updated", zap.Uint("id", reservation.ID))
	return nil
}

func (r *reservationRepository) UpdateStatus(ctx context.Context, id uint, status entity.ReservationStatus) error {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Updating reservation status",
		zap.Uint("id", id),
		zap.String("status", string(status)))

	err := db.Model(&entity.Reservation{}).
		Where("id = ?", id).
		Update("status", status).Error

	if err != nil {
		r.logger.Error("Failed to update reservation status",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Reservation status updated",
		zap.Uint("id", id),
		zap.String("status", string(status)))

	return nil
}

func (r *reservationRepository) Delete(ctx context.Context, id uint) error {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Deleting reservation", zap.Uint("id", id))

	err := db.Delete(&entity.Reservation{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete reservation",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Reservation deleted", zap.Uint("id", id))
	return nil
}
