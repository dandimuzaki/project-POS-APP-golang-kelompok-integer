package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo        UserRepository
	CustomerRepo    CustomerRepository
	TableRepo       TableRepository
	ReservationRepo ReservationRepository
}

func NewRepository(db *gorm.DB, log *zap.Logger) *Repository {
	return &Repository{
		UserRepo:        NewUserRepo(db, log),
		CustomerRepo:    NewCustomerRepository(db, log),
		TableRepo:       NewTableRepository(db, log),
		ReservationRepo: NewReservationRepository(db, log),
	}
}
