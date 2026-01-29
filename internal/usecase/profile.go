package usecase

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/repository"

	"go.uber.org/zap"
)

type ProfileService interface {
	GetProfile(ctx context.Context) (*response.ProfileResponse, error)
	UpdateProfile(ctx context.Context, data *request.ProfileRequest) error
}

type profileService struct {
	tx   TxManager
	repo *repository.Repository
	log  *zap.Logger
}

func NewProfileService(tx TxManager, repo *repository.Repository, log *zap.Logger) ProfileService {
	return &profileService{
		tx: tx,
		repo: repo,
		log: log,
	}
}

func (s *profileService) GetProfile(ctx context.Context) (*response.ProfileResponse, error) {

}

func (s *profileService) UpdateProfile(ctx context.Context, data *request.ProfileRequest) error {
	
}