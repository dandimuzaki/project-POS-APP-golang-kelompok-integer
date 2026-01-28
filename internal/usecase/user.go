package usecase

import (
	"context"
	"math"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto"

	"go.uber.org/zap"
)

type UserService interface{
	GetUserList(ctx context.Context, req dto.UserFilterRequest) ([]dto.UserResponse, dto.Pagination, error)
	GetByID(ctx context.Context, id uint) (entity.User, error)
}

type userService struct {
	repo *repository.Repository
	Logger *zap.Logger
}

func NewUserService(repo *repository.Repository, log *zap.Logger) UserService {
	return &userService{
		repo: repo,
		Logger: log,
	}
}

func (s *userService) GetUserList(ctx context.Context, req dto.UserFilterRequest) ([]dto.UserResponse, dto.Pagination, error) {
	users, total, err := s.repo.UserRepo.GetUserList(ctx, req)
	if err != nil {
		s.Logger.Error("Error get users service", zap.Error(err))
		return nil, dto.Pagination{}, err
	}
	
	var res []dto.UserResponse

	for _, u := range users {
		user := dto.UserResponse{
			Email: u.Email,
			Role: u.Role,
		}
		res = append(res, user)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	pagination := dto.Pagination{
		CurrentPage:  req.Page,
		Limit:        req.Limit,
		TotalPages:   totalPages,
		TotalRecords: total,
	}

	return res, pagination, err
}

func (s *userService) GetByID(ctx context.Context, id uint) (entity.User, error) {
	user, err := s.repo.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		s.Logger.Error("Error get user by id service", zap.Error(err))
		return entity.User{}, err
	}
	return user, nil
}