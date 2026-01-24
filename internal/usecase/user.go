package usecase

import (
	"context"
	"math"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto"

	"go.uber.org/zap"
)

type UserService interface{
	GetListUsers(ctx context.Context, req dto.UserFilterRequest) ([]dto.UserResponse, dto.Pagination, error)
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

func (s *userService) GetListUsers(ctx context.Context, req dto.UserFilterRequest) ([]dto.UserResponse, dto.Pagination, error) {
	users, total, err := s.repo.UserRepo.GetListUsers(ctx, req)
	if err != nil {
		s.Logger.Error("Error get users service", zap.Error(err))
		return nil, dto.Pagination{}, err
	}
	
	var res []dto.UserResponse

	for _, u := range users {
		user := dto.UserResponse{
			Name: u.Name,
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