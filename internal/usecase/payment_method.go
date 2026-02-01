package usecase

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"

	"go.uber.org/zap"
)

type PaymentMethodService interface{
	GetAll(ctx context.Context) ([]entity.PaymentMethod, error)
}

type paymentMethodService struct {
	tx   TxManager
	repo *repository.Repository
	log  *zap.Logger
}

func NewPaymentMethodService(tx TxManager, repo *repository.Repository, log *zap.Logger) PaymentMethodService {
	return &paymentMethodService{
		tx:   tx,
		repo: repo,
		log:  log,
	}
}

func (s *paymentMethodService) GetAll(ctx context.Context) ([]entity.PaymentMethod, error) {
	pMethods, err := s.repo.PaymentMethodRepo.GetAll(ctx)
	if err != nil {
		s.log.Error("Error get payment methods", zap.Error(err))
		return nil, err
	}
	return pMethods, nil
}