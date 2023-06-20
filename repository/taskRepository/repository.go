package taskRepository

import (
	"go-todolist-aws/model"
	"go-todolist-aws/request/taskRequest"
	"go-todolist-aws/utils/paginator"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task model.Task) (model.Task, error)
	GetTaskList(request taskRequest.TaskGetListRequest) paginator.Page[model.Task]
	GetTask(id int64) (model.Task, error)
	UpdateTask(task model.Task) (model.Task, error)
	DeleteTask(id int64) error
	FindByTitle(userID int64, title string) (model.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}
