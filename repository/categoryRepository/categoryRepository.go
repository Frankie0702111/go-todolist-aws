package categoryRepository

import (
	"go-todolist-aws/model"
	"go-todolist-aws/utils/paginator"

	"gorm.io/gorm/clause"
)

func (r *categoryRepository) CreateCategory(category model.Category) (model.Category, error) {
	create := r.db.Save(&category)
	if create.Error != nil {
		return category, create.Error
	}

	return category, nil
}

func (r *categoryRepository) GetCategoryList(id int64, name string, page int64, limit int64) paginator.Page[model.Category] {
	var categories []*model.Category
	query := r.db.Model(&categories).Preload(clause.Associations)

	if id > 0 {
		query.Where("id = ?", id)
	}

	if len(name) > 0 {
		query.Where("name like ?", name+"%")
	}

	p := paginator.Page[model.Category]{CurrentPage: page, PageLimit: limit}
	p.SelectPages(query)

	return p
}

func (r *categoryRepository) GetCategory(id int64) (model.Category, error) {
	var category model.Category
	res := r.db.First(&category, "id=?", id)
	if res.Error != nil {
		return category, res.Error
	}

	return category, nil
}

func (r *categoryRepository) FindByName(name string) (model.Category, error) {
	var category model.Category
	if res := r.db.Where("name = ?", name).Take(&category); res.Error != nil {
		return category, res.Error
	}

	return category, nil
}
