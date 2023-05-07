package category

import (
	"go-todolist-aws/repository/categoryRepository"
	categoryReqyest "go-todolist-aws/request/categoryRequest"
	"go-todolist-aws/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController interface {
	Get(c *gin.Context)
}

type categoryController struct {
	// CategoryService categoryService.CategoryService
	CategoryRepository categoryRepository.CategoryRepository
}

func New(db *gorm.DB) CategoryController {
	return &categoryController{
		// CategoryService: categoryService.New(db),
		CategoryRepository: categoryRepository.New(db),
	}
}

func (c *categoryController) Get(ctx *gin.Context) {
	input := &categoryReqyest.CategoryGetRequest{}
	if err := ctx.ShouldBindUri(input); err != nil {
		response := response.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", response.IdInvalid, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	category, categoryErr := c.CategoryRepository.GetCategory(input.Id)
	if categoryErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", categoryErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	if category.ID == 0 {
		response := response.SuccessResponse(http.StatusOK, "Record not found", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Successfully get category", category)
	ctx.JSON(http.StatusOK, response)
	return
}
