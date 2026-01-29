package usecase

import (
	"context"
	"math"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
)

type UserService interface{
	GetUserList(ctx context.Context, req dto.UserFilterRequest) ([]dto.UserResponse, dto.Pagination, error)
	GetByID(ctx context.Context, id uint) (entity.User, error)
}

type userService struct {
	tx TxManager
	repo *repository.Repository
	log *zap.Logger
}

func NewUserService(tx TxManager, repo *repository.Repository, log *zap.Logger) UserService {
	return &userService{
		tx: tx,
		repo: repo,
		log: log,
	}
}

func (s *userService) GetUserList(ctx context.Context, req dto.UserFilterRequest) ([]dto.UserResponse, dto.Pagination, error) {
	users, total, err := s.repo.UserRepo.GetUserList(ctx, req)
	if err != nil {
		s.log.Error("Error get users service", zap.Error(err))
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
		s.log.Error("Error get user by id service", zap.Error(err))
		return entity.User{}, err
	}
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, req request.UserRequest) (*response.UserResponse, error) {
	password, _ := utils.GenerateRandomString(10)
	user := entity.User{
		Email: req.Email,
		Role: entity.UserRole(req.Role),
		PasswordHash: utils.HashPassword(password),
	}
	var createdUser *entity.User
	var profile entity.Profile
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		result, e := s.repo.UserRepo.CreateUser(ctx, &user)
		if e != nil {
			return e
		}
		createdUser = result

		_, e = s.repo.ProfileRepo.CreateProfile(ctx, &profile)
		if e != nil {
			return e
		}

		return e
	})

	if err != nil {
		s.log.Error("Error create user transaction", zap.Error(err))
		return nil, nil
	}

	res := response.UserResponse{
		Email: createdUser.Email,
		Role: createdUser.Email,
		Password: password,
	}
	return &res, err
}