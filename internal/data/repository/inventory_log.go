package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InventoryLogRepository interface {
	GetInventoryLogs(ctx context.Context, f InventoryLogParams) ([]entity.InventoryLog, int64, error)
	CreateInventoryLog(ctx context.Context, inventory *entity.InventoryLog) (*entity.InventoryLog, error)
}

type inventoryLogRepository struct {
	db *gorm.DB
	Logger *zap.Logger
}

type InventoryLogParams struct {
	Offset int
	Limit int
}

func NewInventoryLogRepo(db *gorm.DB, log *zap.Logger) InventoryLogRepository {
	return &inventoryLogRepository{
		db: db,
		Logger: log,
	}
}

func (r *inventoryLogRepository) GetInventoryLogs(ctx context.Context, f InventoryLogParams) ([]entity.InventoryLog, int64, error) {
	db := infra.GetDB(ctx, r.db)
	var logs []entity.InventoryLog
	var total int64
	query := db.Model(&entity.InventoryLog{})
	
	// Get total inventory logs
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(f.Limit).Offset(f.Offset).Find(&logs).Error
	if err != nil {
		r.Logger.Error("Error query get inventory logs", zap.Error(err))
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *inventoryLogRepository) CreateInventoryLog(ctx context.Context, inventory *entity.InventoryLog) (*entity.InventoryLog, error) {
	db := infra.GetDB(ctx, r.db)
	err := db.Create(&inventory).Error
	if err != nil {
		r.Logger.Error("Error query create inventory", zap.Error(err))
		return nil, err
	}
	return inventory, err
}