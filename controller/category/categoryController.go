package category

import (
	"go-todolist-aws/repository/categoryRepository"
	"go-todolist-aws/request/categoryRequest"
	"go-todolist-aws/service/categoryService"
	"go-todolist-aws/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController interface {
	Create(ctx *gin.Context)
	GetByList(ctx *gin.Context)
	Get(c *gin.Context)
}

type categoryController struct {
	CategoryService    categoryService.CategoryService
	CategoryRepository categoryRepository.CategoryRepository
}

func New(db *gorm.DB) CategoryController {
	return &categoryController{
		CategoryService:    categoryService.New(db),
		CategoryRepository: categoryRepository.New(db),
	}
}

func (c *categoryController) Create(ctx *gin.Context) {
	input := categoryRequest.CategoryCreateOrUpdateRequest{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	createCategory, createCategoryErr := c.CategoryService.CreateCategory(input)
	if createCategoryErr != nil {
		if createCategoryErr.Error() == response.Messages[response.DuplicateCreatedData] {
			response := response.ErrorsResponseByCode(response.DuplicateCreatedData, "Failed to process request", response.DuplicateCreatedData, nil)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", createCategoryErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Create Success", createCategory)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *categoryController) GetByList(ctx *gin.Context) {
	input := &categoryRequest.CategoryGetListRequest{}
	if err := ctx.ShouldBindQuery(input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	category := c.CategoryRepository.GetCategoryList(input.Id, input.Name, input.Page, input.Limit)
	response := response.SuccessPageResponse(http.StatusOK, "Successfully get category list", category.CurrentPage, category.PageLimit, category.Total, category.Pages, category.Data)
	ctx.JSON(http.StatusOK, response)
	return
}

func (c *categoryController) Get(ctx *gin.Context) {
	input := &categoryRequest.CategoryGetRequest{}
	if err := ctx.ShouldBindUri(input); err != nil {
		response := response.ErrorsResponseByCode(response.IdInvalid, "Failed to process request", response.IdInvalid, nil)
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
