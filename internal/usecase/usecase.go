package usecase

import (
	"project-POS-APP-golang-integer/internal/data/repository"

	"go.uber.org/zap"
)

type Usecase struct{
	UserService UserService
}

func NewUsecase(repo *repository.Repository, log *zap.Logger) *Usecase {
	return &Usecase{
		UserService: NewUserService(repo, log),
	}
}