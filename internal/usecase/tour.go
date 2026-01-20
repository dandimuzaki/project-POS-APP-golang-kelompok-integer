package usecase

import (
	"context"
	"math"
	"travel-api/internal/data/repository"
	"travel-api/internal/dto"

	"go.uber.org/zap"
)

type TourService interface{
	GetListTours(ctx context.Context, req dto.TourFilterRequest) ([]dto.TourResponse, dto.Pagination, error)
	ScheduleByID(ctx context.Context, id uint) (*dto.ScheduleResponse, error)
}

type tourService struct {
	repo *repository.Repository
	Logger *zap.Logger
}

func NewTourService(repo *repository.Repository, log *zap.Logger) TourService {
	return &tourService{
		repo: repo,
		Logger: log,
	}
}

func (s *tourService) GetListTours(ctx context.Context, req dto.TourFilterRequest) ([]dto.TourResponse, dto.Pagination, error) {
	// Set default limit jika kosong
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 5
	}

	// 1. Panggil Repository
	tours, total, err := s.repo.TourRepo.FindAll(ctx, req)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	var tourResponses []dto.TourResponse

	for _, t := range tours {
		var sumRating int
		for _, r := range t.Reviews {
			sumRating += r.Rating
		}
		avgRating := 0.0
		if len(t.Reviews) > 0 {
			avgRating = float64(sumRating) / float64(len(t.Reviews))
			avgRating = math.Round(avgRating*10) / 10
		}

		imgUrl := ""
		if len(t.Images) > 0 {
			imgUrl = t.Images[0].ImageURL
		}
		remaining := 0
		if len(t.Schedules) > 0 {
			remaining = t.Schedules[0].Quota - t.Schedules[0].BookedCount
		}

		res := dto.TourResponse{
			ID:            t.ID,
			Name:          t.Name,
			Location:      t.Location.City + ", " + t.Location.Country,
			ImageURL:      imgUrl,
			Price:         t.BasePrice,
			Duration:      t.DurationDay,
			Rating:        avgRating,
			ReviewCount:   len(t.Reviews),
			RemainingSeat: remaining,
		}
		tourResponses = append(tourResponses, res)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	pagination := dto.Pagination{
		CurrentPage:  req.Page,
		Limit:        req.Limit,
		TotalPages:   totalPages,
		TotalRecords: total,
	}

	return tourResponses, pagination, nil
}

func (s *tourService) ScheduleByID(ctx context.Context, id uint) (*dto.ScheduleResponse, error) {
	t, err := s.repo.TourRepo.ScheduleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get rating
	var sumRating int
	for _, r := range t.Reviews {
		sumRating += r.Rating
	}
	avgRating := 0.0
	if len(t.Reviews) > 0 {
		avgRating = float64(sumRating) / float64(len(t.Reviews))
		avgRating = math.Round(avgRating*10) / 10
	}

	var res dto.ScheduleResponse
	res.ID = t.Schedules[0].ID
	res.Tour.Name = t.Name
	res.Tour.Description = t.Description
	res.Tour.Rating = avgRating
	res.Tour.ReviewCount = len(t.Reviews)
	res.Destination = t.Location.City + ", " + t.Location.Country
	res.StartDate = t.Schedules[0].StartDate.Format("January 02, 2006")
	res.EndDate = t.Schedules[0].EndDate.Format("January 02, 2006")
	res.Price.BasePrice = t.BasePrice
	res.Price.PriceOverride = t.Schedules[0].PriceOverride
	var images []string
	for _, img := range t.Images {
		images = append(images, img.ImageURL)
	}
	res.Images = images
	return &res, err
}