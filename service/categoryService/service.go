package categoryService

import (
	"go-todolist-aws/model"
	"go-todolist-aws/repository/categoryRepository"
	"go-todolist-aws/request/categoryRequest"

	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(category categoryRequest.CategoryCreateOrUpdateRequest) (model.Category, error)
	UpdateCategory(category categoryRequest.CategoryCreateOrUpdateRequest, id int64) (model.Category, error)
}

type categoryService struct {
	CategoryRepository categoryRepository.CategoryRepository
}

func New(db *gorm.DB) CategoryService {
	return &categoryService{
		CategoryRepository: categoryRepository.New(db),
	}
}
