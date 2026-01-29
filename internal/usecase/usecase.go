package usecase

import (
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Usecase struct {
	UserService        UserService
	AuthService        AuthService
	ReservationService ReservationService
}

func NewUsecase(tx TxManager, db *gorm.DB, repo *repository.Repository, log *zap.Logger, config utils.Configuration) *Usecase {
	return &Usecase{
		UserService:        NewUserService(tx, repo, log),
		AuthService:        NewAuthService(tx, repo, log, config),
		ReservationService: NewReservationService(tx, repo, log),
	}
}
