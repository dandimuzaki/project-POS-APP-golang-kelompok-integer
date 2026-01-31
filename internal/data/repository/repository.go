package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo        UserRepository
	ProfileRepo ProfileRepository
	OTPRepo OTPRepository
	SessionRepo SessionRepository
	PasswordResetRepo PasswordResetRepository
	CustomerRepo    CustomerRepository
	CategoryRepo CategoryRepository
	TableRepo       TableRepository
	ReservationRepo ReservationRepository
	InventoryLogRepo InventoryLogRepository
}

func NewRepository(db *gorm.DB, log *zap.Logger) *Repository {
	return &Repository{
		UserRepo:        NewUserRepo(db, log),
		ProfileRepo: NewProfileRepo(db, log),
		OTPRepo: NewOTPRepo(db, log),
		SessionRepo: NewSessionRepo(db, log),
		PasswordResetRepo: NewPasswordResetRepo(db, log),
		CategoryRepo: NewCategoryRepository(db),
		CustomerRepo:    NewCustomerRepo(db, log),
		TableRepo:       NewTableRepo(db, log),
		ReservationRepo: NewReservationRepo(db, log),
		InventoryLogRepo: NewInventoryLogRepo(db, log),
	}
}
