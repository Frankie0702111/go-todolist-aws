package taskService

import (
	"errors"
	"go-todolist-aws/model"
	"go-todolist-aws/request/taskRequest"
	"go-todolist-aws/utils/aws"
	"go-todolist-aws/utils/response"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofrs/uuid"
	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

func initS3() *s3.Client {
	s3, s3Err := aws.InitS3()
	if s3Err != nil {
		return nil
	}

	return s3
}

func (s *taskService) CreateTask(task taskRequest.TaskCreateRequest) (model.Task, error) {
	taskToCteate := model.Task{}
	if err := smapping.FillStruct(&taskToCteate, smapping.MapFields(&task)); err != nil {
		return taskToCteate, err
	}

	_, err := s.TaskRepository.FindByTitle(task.UserID, task.Title)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if task.Image != nil {
				uuidV4Ojb, uuidV4Err := uuid.NewV4()
				if uuidV4Err != nil {
					return taskToCteate, uuidV4Err
				}

				s3Upload, err := s.S3Repository.FileUpload(task.Image, uuidV4Ojb.String())
				if err != nil {
					return taskToCteate, err
				}

				taskToCteate.Img = task.Image.Filename
				taskToCteate.ImgLink = s3Upload.Location
				taskToCteate.ImgUuid = uuidV4Ojb.String()
			}

			res, err := s.TaskRepository.CreateTask(taskToCteate)
			if err != nil {
				return res, err
			}

			return res, nil
		} else {
			return taskToCteate, err
		}
	}

	return taskToCteate, errors.New(response.Messages[response.DuplicateCreatedData])
}

func (s *taskService) UpdateTask(task taskRequest.TaskUpdateRequest, old_task model.Task) (model.Task, error) {
	taskToUpdate := model.Task{}
	if err := smapping.FillStruct(&taskToUpdate, smapping.MapFields(&task)); err != nil {
		return taskToUpdate, err
	}

	checkTitle, _ := s.TaskRepository.FindByTitle(old_task.UserID, task.Title)
	if (task.Title != checkTitle.Title) || ((old_task.ID == checkTitle.ID) && (task.Title == checkTitle.Title)) {
		if task.Image != nil {
			if task.Image.Filename != old_task.Img {
				if err := s.S3Repository.FileRemove(old_task.Img, old_task.ImgUuid); err != nil {
					return taskToUpdate, err
				}
			}

			uuidV4Ojb, uuidV4Err := uuid.NewV4()
			if uuidV4Err != nil {
				return taskToUpdate, uuidV4Err
			}

			s3Upload, err := s.S3Repository.FileUpload(task.Image, uuidV4Ojb.String())
			if err != nil {
				return taskToUpdate, err
			}

			taskToUpdate.Img = task.Image.Filename
			taskToUpdate.ImgLink = s3Upload.Location
			taskToUpdate.ImgUuid = uuidV4Ojb.String()
		}

		taskToUpdate.ID = old_task.ID
		res, err := s.TaskRepository.UpdateTask(taskToUpdate)
		if err != nil {
			return res, err
		}

		return res, nil
	}

	return taskToUpdate, errors.New(response.Messages[response.DuplicatedTitle])
}

func (s *taskService) DeleteTask(task model.Task) error {
	if err := s.S3Repository.FileRemove(task.Img, task.ImgUuid); err != nil {
		return err
	}

	if err := s.TaskRepository.DeleteTask(task.ID); err != nil {
		return err
	}

	return nil
}
