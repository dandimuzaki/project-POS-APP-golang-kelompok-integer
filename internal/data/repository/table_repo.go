package repository

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TableRepository interface {
	Create(table *entity.Table, tx ...*gorm.DB) error
	FindByID(id uint, tx ...*gorm.DB) (*entity.Table, error)
	FindByNumber(tableNumber string, tx ...*gorm.DB) (*entity.Table, error)
	FindAll(params TableQueryParams, tx ...*gorm.DB) ([]entity.Table, int64, error)
	FindByCapacity(minCapacity int, tx ...*gorm.DB) ([]entity.Table, error)
	FindByStatus(status entity.TableStatus, tx ...*gorm.DB) ([]entity.Table, error)
	Update(table *entity.Table, tx ...*gorm.DB) error
	UpdateStatus(id uint, status entity.TableStatus, tx ...*gorm.DB) error
	Delete(id uint, tx ...*gorm.DB) error
}

type tableRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewTableRepository(db *gorm.DB, log *zap.Logger) TableRepository {
	return &tableRepository{
		db:  db,
		log: log.With(zap.String("repository", "table")),
	}
}

type TableQueryParams struct {
	Status      entity.TableStatus // Filter by status
	MinCapacity int                // Minimal capacity
	Offset      int                // Pagination offset
	Limit       int                // Pagination limit
}

// getDB returns either transaction db or regular db
func (tr *tableRepository) getDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}
	return tr.db
}

func (tr *tableRepository) Create(table *entity.Table, tx ...*gorm.DB) error {
	db := tr.getDB(tx...)

	tr.log.Debug("Creating new table",
		zap.String("table_number", table.TableNumber))

	err := db.Create(table).Error
	if err != nil {
		tr.log.Error("Failed to create table",
			zap.String("table_number", table.TableNumber),
			zap.Error(err))
		return err
	}

	tr.log.Info("Table created successfully",
		zap.Uint("id", table.ID),
		zap.String("table_number", table.TableNumber))
	return nil
}

func (tr *tableRepository) FindByID(id uint, tx ...*gorm.DB) (*entity.Table, error) {
	db := tr.getDB(tx...)

	tr.log.Debug("Finding table by ID", zap.Uint("id", id))

	var table entity.Table
	err := db.First(&table, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tr.log.Warn("Table not found", zap.Uint("id", id))
		} else {
			tr.log.Error("Failed to find table by ID",
				zap.Uint("id", id),
				zap.Error(err))
		}
		return nil, err
	}

	return &table, nil
}

func (tr *tableRepository) FindByNumber(tableNumber string, tx ...*gorm.DB) (*entity.Table, error) {
	db := tr.getDB(tx...)

	tr.log.Debug("Finding table by number", zap.String("table_number", tableNumber))

	var table entity.Table
	err := db.Where("table_number = ?", tableNumber).First(&table).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tr.log.Debug("Table not found by number", zap.String("table_number", tableNumber))
		} else {
			tr.log.Error("Failed to find table by number",
				zap.String("table_number", tableNumber),
				zap.Error(err))
		}
		return nil, err
	}

	return &table, nil
}

func (tr *tableRepository) FindAll(params TableQueryParams, tx ...*gorm.DB) ([]entity.Table, int64, error) {
	db := tr.getDB(tx...)

	tr.log.Debug("Finding tables",
		zap.String("status", string(params.Status)),
		zap.Int("min_capacity", params.MinCapacity),
		zap.Int("offset", params.Offset),
		zap.Int("limit", params.Limit))

	var tables []entity.Table
	var total int64

	query := db.Model(&entity.Table{})

	// Apply status filter if provided
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Apply capacity filter if provided
	if params.MinCapacity > 0 {
		query = query.Where("capacity >= ?", params.MinCapacity)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		tr.log.Error("Failed to count tables", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Offset(params.Offset).Limit(params.Limit)
	}

	// Execute query with ordering
	if err := query.Order("table_number ASC").Find(&tables).Error; err != nil {
		tr.log.Error("Failed to find tables",
			zap.Error(err))
		return nil, 0, err
	}

	tr.log.Debug("Tables found",
		zap.Int("count", len(tables)),
		zap.Int64("total", total))
	return tables, total, nil
}

func (tr *tableRepository) FindByCapacity(minCapacity int, tx ...*gorm.DB) ([]entity.Table, error) {
	db := tr.getDB(tx...)

	tr.log.Debug("Finding tables by capacity", zap.Int("min_capacity", minCapacity))

	var tables []entity.Table
	err := db.
		Where("capacity >= ?", minCapacity).
		Order("capacity ASC"). // Prioritize smaller tables first
		Find(&tables).Error

	if err != nil {
		tr.log.Error("Failed to find tables by capacity",
			zap.Int("min_capacity", minCapacity),
			zap.Error(err))
		return nil, err
	}

	tr.log.Debug("Tables found by capacity",
		zap.Int("min_capacity", minCapacity),
		zap.Int("count", len(tables)))
	return tables, nil
}

func (tr *tableRepository) FindByStatus(status entity.TableStatus, tx ...*gorm.DB) ([]entity.Table, error) {
	db := tr.getDB(tx...)

	tr.log.Debug("Finding tables by status", zap.String("status", string(status)))

	var tables []entity.Table
	err := db.
		Where("status = ?", status).
		Order("table_number ASC").
		Find(&tables).Error

	if err != nil {
		tr.log.Error("Failed to find tables by status",
			zap.String("status", string(status)),
			zap.Error(err))
		return nil, err
	}

	return tables, nil
}

func (tr *tableRepository) Update(table *entity.Table, tx ...*gorm.DB) error {
	db := tr.getDB(tx...)

	tr.log.Debug("Updating table",
		zap.Uint("id", table.ID),
		zap.String("table_number", table.TableNumber))

	err := db.Save(table).Error
	if err != nil {
		tr.log.Error("Failed to update table",
			zap.Uint("id", table.ID),
			zap.Error(err))
		return err
	}

	tr.log.Info("Table updated", zap.Uint("id", table.ID))
	return nil
}

func (tr *tableRepository) UpdateStatus(id uint, status entity.TableStatus, tx ...*gorm.DB) error {
	db := tr.getDB(tx...)

	tr.log.Debug("Updating table status",
		zap.Uint("id", id),
		zap.String("status", string(status)))

	err := db.Model(&entity.Table{}).
		Where("id = ?", id).
		Update("status", status).Error

	if err != nil {
		tr.log.Error("Failed to update table status",
			zap.Uint("id", id),
			zap.String("status", string(status)),
			zap.Error(err))
		return err
	}

	tr.log.Info("Table status updated",
		zap.Uint("id", id),
		zap.String("status", string(status)))
	return nil
}

func (tr *tableRepository) Delete(id uint, tx ...*gorm.DB) error {
	db := tr.getDB(tx...)

	tr.log.Debug("Deleting table", zap.Uint("id", id))

	err := db.Delete(&entity.Table{}, id).Error
	if err != nil {
		tr.log.Error("Failed to delete table",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	tr.log.Info("Table deleted", zap.Uint("id", id))
	return nil
}
