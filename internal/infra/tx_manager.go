package infra

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

type GormTxManager struct {
	db *gorm.DB
}

func NewGormTxManager(db *gorm.DB) *GormTxManager {
	return &GormTxManager{db: db}
}

func (tm *GormTxManager) WithinTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}
