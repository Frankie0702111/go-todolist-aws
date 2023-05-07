package categoryRepository

import (
	"go-todolist-aws/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategory(id int64) (res model.Category, err error)
}

type categoryRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
