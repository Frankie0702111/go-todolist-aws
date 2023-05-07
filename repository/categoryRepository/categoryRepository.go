package categoryRepository

import "go-todolist-aws/model"

func (r *categoryRepository) GetCategory(id int64) (category model.Category, err error) {
	res := r.db.First(&category, "id=?", id)
	if res.Error == nil {
		return category, nil
	}

	return category, err
}
