package categoryService

import (
	"errors"
	"go-todolist-aws/model"
	"go-todolist-aws/request/categoryRequest"
	"go-todolist-aws/utils/response"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

func (s *categoryService) CreateCategory(category categoryRequest.CategoryCreateOrUpdateRequest) (model.Category, error) {
	createCategory := model.Category{}
	if err := smapping.FillStruct(&createCategory, smapping.MapFields(&category)); err != nil {
		return createCategory, err
	}

	_, err := s.CategoryRepository.FindByName(category.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res, err := s.CategoryRepository.CreateCategory(createCategory)
			if err != nil {
				return res, err
			}

			return res, nil
		} else {
			return createCategory, err
		}
	}

	return createCategory, errors.New(response.Messages[response.DuplicateCreatedData])
}

func (s *categoryService) UpdateCategory(category categoryRequest.CategoryCreateOrUpdateRequest, id int64) (model.Category, error) {
	categoryToUpdate := model.Category{}
	if err := smapping.FillStruct(&categoryToUpdate, smapping.MapFields(&category)); err != nil {
		return categoryToUpdate, err
	}

	categoryToUpdate.ID = id
	res, err := s.CategoryRepository.UpdateCategory(categoryToUpdate)
	if err != nil {
		return res, err
	}

	return res, nil

}
