package usecase

import (
	"project-POS-APP-golang-integer/internal/data/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Usecase struct {
	UserService        UserService
	ReservationService ReservationService
}

func NewUsecase(db *gorm.DB, repo *repository.Repository, log *zap.Logger) *Usecase {
	return &Usecase{
		UserService:        NewUserService(repo, log),
		ReservationService: NewReservationService(db, repo, log),
	}
}
