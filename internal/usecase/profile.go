package usecase

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"time"

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
	userID := ctx.Value("user_id").(uint)
	profile, err := s.repo.ProfileRepo.GetProfileByID(ctx, userID)
	if err != nil {
		s.log.Error("Error get profile", zap.Error(err))
		return nil, err
	}
	var birthday string
	if profile.Profile != nil && profile.Profile.DateOfBirth != nil {
		birthday = profile.Profile.DateOfBirth.Format("2 January 2006")
	}

	res := response.ProfileResponse{
		Email: profile.Email,
		Role:  profile.Role,
	}

	if profile.Profile != nil {
		res.FullName = profile.Profile.FullName
		res.Phone = profile.Profile.Phone
		res.DateOfBirth = birthday
		res.Salary = profile.Profile.Salary
		res.ProfileImageURL = profile.Profile.ProfileImageURL
		res.Address = profile.Profile.Address
		res.AdditionalDetails = profile.Profile.AdditionalDetails
	}

	return &res, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, data *request.ProfileRequest) error {
	userID := ctx.Value("user_id").(uint)
	birthday, _ := time.Parse("2 January 2006", data.DateOfBirth)
	profile := entity.Profile{
		UserID: userID,
		FullName: data.FullName,
		Phone: data.Phone,
		DateOfBirth: &birthday,
		Salary: data.Salary,
		ProfileImageURL: data.ProfileImageURL,
		Address: data.Address,
		AdditionalDetails: data.AdditionalDetails,
	}
	err := s.repo.ProfileRepo.UpdateProfile(ctx, &profile)
	if err != nil {
		s.log.Error("Error update profile service", zap.Error(err))
		return err
	}
	return nil
}