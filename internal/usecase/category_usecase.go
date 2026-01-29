package usecase

import (
	"context"

	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto"
)

type CategoryUsecase interface {
	GetAllCategories(ctx context.Context) ([]entity.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error)
	CreateCategory(ctx context.Context, input dto.CreateCategoryInput) (*entity.Category, error)
	UpdateCategory(ctx context.Context, id int64, input dto.UpdateCategoryInput) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (u *categoryUsecase) GetAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := u.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return mapToCategoryResponses(categories), nil
}

// ===== mapper =====

func mapToCategoryResponses(categories []entity.Category) []dto.CategoryResponse {
	res := make([]dto.CategoryResponse, 0)
	for _, c := range categories {
		res = append(res, dto.CategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}
	return res
}
