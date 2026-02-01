package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PaymentMethodRepository interface{
	GetAll(ctx context.Context) ([]entity.PaymentMethod, error)
}

type paymentMethodRepository struct {
	db     *gorm.DB
	Logger *zap.Logger
}

func NewPaymentMethodRepo(db *gorm.DB, log *zap.Logger) PaymentMethodRepository {
	return &paymentMethodRepository{
		db:     db,
		Logger: log,
	}
}

func (r *paymentMethodRepository) GetAll(ctx context.Context) ([]entity.PaymentMethod, error) {
	db := infra.GetDB(ctx, r.db)
	var pMethods []entity.PaymentMethod
	err := db.Model(&entity.PaymentMethod{}).Find(&pMethods).Error
	if err != nil {
		r.Logger.Error("Error query get payment methods", zap.Error(err))
		return nil, err
	}

	return pMethods, nil
}
