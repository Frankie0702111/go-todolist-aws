package categoryRepository

import (
	"go-todolist-aws/model"
	"go-todolist-aws/request/categoryRequest"
	"go-todolist-aws/utils/paginator"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category model.Category) (model.Category, error)
	GetCategoryList(request categoryRequest.CategoryGetListRequest) paginator.Page[model.Category]
	GetCategory(id int64) (model.Category, error)
	UpdateCategory(category model.Category) (model.Category, error)
	DeleteCategory(id int64) error
	FindByName(name string) (model.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
