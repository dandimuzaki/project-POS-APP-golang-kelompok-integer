package usecase

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/pkg/utils"

	"go.uber.org/zap"
)

type UserService interface{
	GetUserList(ctx context.Context, req request.UserFilterRequest) (*response.PaginatedResponse[response.UserResponse], error)
	GetUserByID(ctx context.Context, id uint) (entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, req request.UserRequest) (*response.CreateUserResponse, error)
	UpdateRole(ctx context.Context, req request.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id uint) error
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

func (s *userService) GetUserList(ctx context.Context, req request.UserFilterRequest) (*response.PaginatedResponse[response.UserResponse], error) {
	filter := repository.UserQueryParams{
		Offset: req.GetOffset(),
		Limit: req.GetPerPage(),
		Role: entity.UserRole(req.Role),
		Name: req.Name,
		Email: req.Email,
	}
	
	users, total, err := s.repo.UserRepo.GetUserList(ctx, filter)
	if err != nil {
		s.log.Error("Error get users service", zap.Error(err))
		return nil, err
	}
	
	// Convert to DTO
	var res []response.UserResponse
	for _, u := range users {
		user := response.UserResponse{
			Email: u.Email,
			Role: u.Role,
		}
		res = append(res, user)
	}

	return response.NewPaginatedResponse(
		res,
		req.GetPage(),
		req.GetPerPage(),
		total,
	), nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (entity.User, error) {
	user, err := s.repo.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		s.log.Error("Error get user by id service", zap.Error(err))
		return entity.User{}, err
	}
	return user, nil
}

func (s *userService) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.repo.UserRepo.FindUserByEmail(ctx, email)
	if err != nil {
		s.log.Error("Error get user by id service", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, req request.UserRequest) (*response.CreateUserResponse, error) {
	password, _ := utils.GenerateRandomString(10)
	user := entity.User{
		Email: req.Email,
		Role: entity.UserRole(req.Role),
		PasswordHash: utils.HashPassword(password),
	}
	var createdUser *entity.User
	var profile entity.Profile
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// Create user
		result, e := s.repo.UserRepo.CreateUser(ctx, &user)
		if e != nil {
			return e
		}
		createdUser = result

		// Automatically create empty profile
		profile.UserID = createdUser.ID
		_, e = s.repo.ProfileRepo.CreateProfile(ctx, &profile)
		if e != nil {
			return e
		}
		
		return nil
	})

	if err != nil {
		s.log.Error("Error create user transaction", zap.Error(err))
		return nil, err
	}

	res := response.CreateUserResponse{
		Email: createdUser.Email,
		Role: createdUser.Role,
		Password: password,
	}
	return &res, err
}

func (s *userService) UpdateRole(ctx context.Context, req request.UpdateUserRequest) error {
	user, err := s.repo.UserRepo.GetUserByID(ctx, req.ID)
	if err != nil {
		s.log.Error("Error get user by id", zap.Error(err))
		return err
	}
	
	// Update user role
	user.Role = entity.UserRole(req.Role)
	err = s.repo.UserRepo.UpdateUser(ctx, req.ID, &user)
	if err != nil {
		s.log.Error("Error update user", zap.Error(err))
		return err
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	err := s.repo.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		s.log.Error("Error delete user", zap.Error(err))
		return err
	}

	return nil
}