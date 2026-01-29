package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TableRepository interface {
	Create(ctx context.Context, table *entity.Table) (*entity.Table, error)
	FindByID(ctx context.Context, id uint) (*entity.Table, error)
	FindByNumber(ctx context.Context, tableNumber string) (*entity.Table, error)
	FindAll(ctx context.Context, params request.GetTablesRequest) ([]entity.Table, int64, error)
	FindByCapacity(ctx context.Context, minCapacity int) ([]entity.Table, error)
	FindByStatus(ctx context.Context, status entity.TableStatus) ([]entity.Table, error)
	Update(ctx context.Context, table *entity.Table) error
	UpdateStatus(ctx context.Context, id uint, status entity.TableStatus) error
	Delete(ctx context.Context, id uint) error
}

type tableRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTableRepo(db *gorm.DB, log *zap.Logger) TableRepository {
	return &tableRepository{
		db:     db,
		logger: log.With(zap.String("repository", "table")),
	}
}

func (r *tableRepository) Create(ctx context.Context, table *entity.Table) (*entity.Table, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Creating table",
		zap.String("table_number", table.TableNumber),
		zap.Int("capacity", table.Capacity))

	err := db.Create(table).Error
	if err != nil {
		r.logger.Error("Failed to create table",
			zap.String("table_number", table.TableNumber),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Table created successfully",
		zap.Uint("id", table.ID),
		zap.String("table_number", table.TableNumber))

	return table, nil
}

func (r *tableRepository) FindByID(ctx context.Context, id uint) (*entity.Table, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding table by ID", zap.Uint("id", id))

	var table entity.Table
	err := db.First(&table, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn("Table not found", zap.Uint("id", id))
		} else {
			r.logger.Error("Failed to find table",
				zap.Uint("id", id),
				zap.Error(err))
		}
		return nil, err
	}

	return &table, nil
}

func (r *tableRepository) FindByNumber(ctx context.Context, tableNumber string) (*entity.Table, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding table by number", zap.String("table_number", tableNumber))

	var table entity.Table
	err := db.Where("table_number = ?", tableNumber).First(&table).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Table not found by number", zap.String("table_number", tableNumber))
		} else {
			r.logger.Error("Failed to find table by number",
				zap.String("table_number", tableNumber),
				zap.Error(err))
		}
		return nil, err
	}

	return &table, nil
}

func (r *tableRepository) FindAll(ctx context.Context, params request.GetTablesRequest) ([]entity.Table, int64, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding tables",
		zap.String("status", params.Status),
		zap.Int("min_capacity", params.MinCapacity),
		zap.Int("page", params.GetPage()),
		zap.Int("per_page", params.GetPerPage()))

	var tables []entity.Table
	var total int64

	query := db.Model(&entity.Table{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.MinCapacity > 0 {
		query = query.Where("capacity >= ?", params.MinCapacity)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count tables", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	offset := params.GetOffset()
	limit := params.GetPerPage()

	err := query.
		Offset(offset).
		Limit(limit).
		Order("table_number ASC").
		Find(&tables).Error

	if err != nil {
		r.logger.Error("Failed to find tables",
			zap.Error(err),
			zap.Int("offset", offset),
			zap.Int("limit", limit))
		return nil, 0, err
	}

	r.logger.Debug("Tables found",
		zap.Int("count", len(tables)),
		zap.Int64("total", total))

	return tables, total, nil
}

func (r *tableRepository) FindByCapacity(ctx context.Context, minCapacity int) ([]entity.Table, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding tables by capacity", zap.Int("min_capacity", minCapacity))

	var tables []entity.Table
	err := db.
		Where("capacity >= ?", minCapacity).
		Order("capacity ASC").
		Find(&tables).Error

	if err != nil {
		r.logger.Error("Failed to find tables by capacity",
			zap.Int("min_capacity", minCapacity),
			zap.Error(err))
		return nil, err
	}

	r.logger.Debug("Tables found by capacity",
		zap.Int("min_capacity", minCapacity),
		zap.Int("count", len(tables)))

	return tables, nil
}

func (r *tableRepository) FindByStatus(ctx context.Context, status entity.TableStatus) ([]entity.Table, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding tables by status", zap.String("status", string(status)))

	var tables []entity.Table
	err := db.
		Where("status = ?", status).
		Order("table_number ASC").
		Find(&tables).Error

	if err != nil {
		r.logger.Error("Failed to find tables by status",
			zap.String("status", string(status)),
			zap.Error(err))
		return nil, err
	}

	return tables, nil
}

func (r *tableRepository) Update(ctx context.Context, table *entity.Table) error {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Updating table",
		zap.Uint("id", table.ID),
		zap.String("table_number", table.TableNumber))

	err := db.Save(table).Error
	if err != nil {
		r.logger.Error("Failed to update table",
			zap.Uint("id", table.ID),
			zap.Error(err))
		return err
	}

	r.logger.Info("Table updated", zap.Uint("id", table.ID))
	return nil
}

func (r *tableRepository) UpdateStatus(ctx context.Context, id uint, status entity.TableStatus) error {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Updating table status",
		zap.Uint("id", id),
		zap.String("status", string(status)))

	err := db.Model(&entity.Table{}).
		Where("id = ?", id).
		Update("status", status).Error

	if err != nil {
		r.logger.Error("Failed to update table status",
			zap.Uint("id", id),
			zap.String("status", string(status)),
			zap.Error(err))
		return err
	}

	r.logger.Info("Table status updated",
		zap.Uint("id", id),
		zap.String("status", string(status)))

	return nil
}

func (r *tableRepository) Delete(ctx context.Context, id uint) error {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Deleting table", zap.Uint("id", id))

	err := db.Delete(&entity.Table{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete table",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Table deleted", zap.Uint("id", id))
	return nil
}
