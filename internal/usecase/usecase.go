package usecase

import (
	"project-POS-APP-golang-integer/internal/data/repository"

	"go.uber.org/zap"
)

type Usecase struct {
	UserService         UserService
	AuthService         AuthService
	CategoryService     CategoryService
	ProfileService      ProfileService
	ReservationService  ReservationService
	InventoryLogService InventoryLogService
	ProductService      ProductService
}

func NewUsecase(tx TxManager, repo *repository.Repository, log *zap.Logger, email EmailSender) *Usecase {
	return &Usecase{
		UserService:         NewUserService(tx, repo, log, email),
		AuthService:         NewAuthService(tx, repo, log, email),
		ProfileService:      NewProfileService(tx, repo, log),
		CategoryService:     NewCategoryService(tx, repo.Category, log),
		ProductService:      NewProductService(tx, repo.Product, repo.Category, log),
		ReservationService:  NewReservationService(tx, repo, log),
		InventoryLogService: NewInventoryLogService(tx, repo, log),
	}
}
