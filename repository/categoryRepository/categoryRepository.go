package categoryRepository

import (
	"go-todolist-aws/model"
	"go-todolist-aws/request/categoryRequest"
	"go-todolist-aws/utils/paginator"
)

func (r *categoryRepository) CreateCategory(category model.Category) (model.Category, error) {
	if create := r.db.Save(&category); create.Error != nil {
		return category, create.Error
	}

	return category, nil
}

func (r *categoryRepository) GetCategoryList(request categoryRequest.CategoryGetListRequest) paginator.Page[model.Category] {
	var categories []*model.Category
	query := r.db.Model(&categories)

	if request.Id > 0 {
		query.Where("id = ?", request.Id)
	}

	if len(request.Name) > 0 {
		query.Where("name like ?", request.Name+"%")
	}

	p := paginator.Page[model.Category]{CurrentPage: request.Page, PageLimit: request.Limit}
	p.SelectPages(query)

	return p
}

func (r *categoryRepository) GetCategory(id int64) (model.Category, error) {
	var category model.Category
	if get := r.db.First(&category, "id=?", id); get.Error != nil {
		return category, get.Error
	}

	return category, nil
}

func (r *categoryRepository) UpdateCategory(category model.Category) (model.Category, error) {
	if update := r.db.Where("id = ?", category.ID).Updates(&category); update.Error != nil {
		return category, update.Error
	}

	return category, nil
}

func (r *categoryRepository) DeleteCategory(id int64) error {
	var category model.Category
	if delete := r.db.Delete(&category, id); delete.Error != nil {
		return delete.Error
	}

	return nil
}

func (r *categoryRepository) FindByName(name string) (model.Category, error) {
	var category model.Category
	if res := r.db.Where("name = ?", name).Take(&category); res.Error != nil {
		return category, res.Error
	}

	return category, nil
}
