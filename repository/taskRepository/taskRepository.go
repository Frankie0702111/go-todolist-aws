package taskRepository

import (
	"go-todolist-aws/model"
	"go-todolist-aws/request/taskRequest"
	"go-todolist-aws/utils/paginator"

	"gorm.io/gorm/clause"
)

func (r *taskRepository) CreateTask(task model.Task) (model.Task, error) {
	if create := r.db.Save(&task); create.Error != nil {
		return task, create.Error
	}

	return task, nil
}

func (r *taskRepository) GetTaskList(request taskRequest.TaskGetListRequest) paginator.Page[model.Task] {
	var task []*model.Task
	query := r.db.Model(&task).Preload(clause.Associations)

	if request.Id > 0 {
		query.Where("id = ?", request.Id)
	}

	if request.UserID > 0 {
		query.Where("id = ?", request.UserID)
	}

	if len(request.Title) > 0 {
		query.Where("title like ?", request.Title+"%")
	}

	if request.SpecifyDatetime != nil {
		query.Where("specify_datetime = ?", request.SpecifyDatetime)
	}

	if request.IsSpecifyTime != nil {
		query.Where("is_specify_time = ?", request.IsSpecifyTime)
	}

	if request.IsComplete != nil {
		query.Where("is_complete = ?", request.IsComplete)
	}

	p := paginator.Page[model.Task]{CurrentPage: request.Page, PageLimit: request.Limit}
	p.SelectPages(query)

	return p
}

func (r *taskRepository) GetTask(id int64) (model.Task, error) {
	var task model.Task
	if get := r.db.Preload("Category").First(&task, "id=?", id); get.Error != nil {
		return task, get.Error
	}

	return task, nil
}

func (r *taskRepository) UpdateTask(task model.Task) (model.Task, error) {
	if update := r.db.Where("id = ?", task.ID).Updates(&task); update.Error != nil {
		return task, update.Error
	}

	return task, nil
}

func (r *taskRepository) DeleteTask(id int64) error {
	var task model.Task
	if delete := r.db.Delete(&task, id); delete.Error != nil {
		return delete.Error
	}

	return nil
}

func (r *taskRepository) FindByTitle(user_id int64, title string) (model.Task, error) {
	var task model.Task
	if res := r.db.Where("user_id = ? AND title = ?", user_id, title).Take(&task); res.Error != nil {
		return task, res.Error
	}

	return task, nil
}
