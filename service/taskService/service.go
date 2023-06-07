package taskService

import (
	"go-todolist-aws/model"
	"go-todolist-aws/repository/s3Repository"
	"go-todolist-aws/repository/taskRepository"
	"go-todolist-aws/request/taskRequest"

	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(task taskRequest.TaskCreateRequest) (model.Task, error)
	UpdateTask(task taskRequest.TaskUpdateRequest, old_task model.Task) (model.Task, error)
	DeleteTask(task model.Task) error
}

type taskService struct {
	TaskRepository taskRepository.TaskRepository
	S3Repository   s3Repository.S3Repository
}

func New(db *gorm.DB) TaskService {
	return &taskService{
		TaskRepository: taskRepository.New(db),
		S3Repository:   s3Repository.New(initS3()),
	}
}
