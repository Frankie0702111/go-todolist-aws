package categoryRepository

import (
	"go-todolist-aws/model"
	"go-todolist-aws/utils/paginator"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category model.Category) (model.Category, error)
	GetCategoryList(id int64, name string, page int64, limit int64) paginator.Page[model.Category]
	GetCategory(id int64) (model.Category, error)
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
