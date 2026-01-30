package usecase

import (
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/response"
)

type CategoryUsecase interface {
	GetAllCategories() ([]response.CategoryResponse, error)
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{categoryRepo}
}

func (uc *categoryUsecase) GetAllCategories() ([]response.CategoryResponse, error) {
	// Get categories from repository
	categories, err := uc.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Map entity to DTO response
	var categoryResponses []response.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, response.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			IconURL:     category.IconURL,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	return categoryResponses, nil
}
