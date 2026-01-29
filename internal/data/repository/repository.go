package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo        UserRepository
	ProfileRepo     ProfileRepository
	SessionRepo     SessionRepository
	CustomerRepo    CustomerRepository
	TableRepo       TableRepository
	ReservationRepo ReservationRepository
}

func NewRepository(db *gorm.DB, log *zap.Logger) *Repository {
	return &Repository{
		UserRepo:        NewUserRepo(db, log),
		ProfileRepo:     NewProfileRepo(db, log),
		SessionRepo:     NewSessionRepo(db, log),
		CustomerRepo:    NewCustomerRepo(db, log),
		TableRepo:       NewTableRepo(db, log),
		ReservationRepo: NewReservationRepo(db, log),
	}
}
