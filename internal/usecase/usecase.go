package usecase

import (
	"travel-api/internal/data/repository"

	"go.uber.org/zap"
)

type Usecase struct{
	TourService TourService
}

func NewUsecase(repo *repository.Repository, log *zap.Logger) *Usecase {
	return &Usecase{
		TourService: NewTourService(repo, log),
	}
}