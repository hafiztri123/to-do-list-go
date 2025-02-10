package usecase

import (
	"github.com/hafiztri123/internal/adapters/secondary/persistent"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/response"
)


type CategoryService struct {
	repo *persistent.CategoryRepository
}

func NewCategoryService(repo *persistent.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(userid uint, name string) error {

	if err := s.repo.Create(userid, name); err != nil {
		return &response.AppError{
			Code:    500,
			Message: err.Error(),
		}
	}
	return nil

}

func (s *CategoryService) GetAllCategory(userID uint) (*[]entity.Category, error) {

	category, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, &response.AppError{
			Code:    500,
			Message: err.Error(),
		}
	}

	return category, nil

}

func (s *CategoryService) DeleteCategory(id uint) error {
	if !s.repo.IsCategoryExistByID(id) {
		return &response.AppError{
			Code:    404,
			Message: "Category not found",
		}
	}

	if err := s.repo.Delete(id); err != nil {
		return &response.AppError{
			Code:    500,
			Message: err.Error(),
		}
	}
	return nil
}