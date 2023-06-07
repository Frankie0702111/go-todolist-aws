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
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
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

// @Summary "Create category"
// @Tags	"Category"
// @Version 1.0
// @Produce application/json
// @Param	Authorization header string true "example:Bearer token (Bearer+space+token)." default(Bearer )
// @Param	* body categoryRequest.CategoryCreateOrUpdateRequest true "Create category"
// @Success 200 object response.Response{errors=string,data=string} "Create Success"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/category [post]
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

// @Summary "Category list"
// @Tags	"Category"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header		string	true	"example:Bearer token (Bearer+space+token)."	default(Bearer )
// @Param	id				query		integer	false	"Category ID"									minimum(1)
// @Param	name			query		string	false	"Category Name"									maxLength(100)
// @Param	page			query		integer	true	"Page"											minimum(1) default(1)
// @Param	limit			query		integer	true	"Limit"											minimum(2) default(5)
// @Success 200 object response.PageResponse{errors=string,data=string} "Successfully get category list"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/category [get]
func (c *categoryController) GetByList(ctx *gin.Context) {
	input := categoryRequest.CategoryGetListRequest{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	category := c.CategoryRepository.GetCategoryList(input)
	response := response.SuccessPageResponse(http.StatusOK, "Successfully get category list", category.CurrentPage, category.PageLimit, category.Total, category.Pages, category.Data)
	ctx.JSON(http.StatusOK, response)
	return
}

// @Summary "Get a single category"
// @Tags	"Category"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header	string	true	"example:Bearer token (Bearer+space+token)."	default(Bearer )
// @Param	id				path	integer	true	"Category ID"									minimum(1)
// @Success 200 object response.Response{errors=string,data=string} "Record not found || Successfully get category"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/category/{id} [get]
func (c *categoryController) Get(ctx *gin.Context) {
	input := &categoryRequest.CategoryGetRequest{}
	if err := ctx.ShouldBindUri(input); err != nil {
		response := response.ErrorsResponseByCode(response.IdInvalid, "Failed to process request", response.IdInvalid, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	category, categoryErr := c.CategoryRepository.GetCategory(input.Id)
	if category.ID == 0 {
		response := response.SuccessResponse(http.StatusOK, "Record not found", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}
	if categoryErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", categoryErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Successfully get category", category)
	ctx.JSON(http.StatusOK, response)
	return
}

// @Summary "Update category"
// @Tags	"Category"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header	string	true	"example:Bearer token (Bearer+space+token)."	default(Bearer )
// @Param	id				path	integer	true	"Category ID"									minimum(1)
// @Param	* body categoryRequest.CategoryCreateOrUpdateRequest true "Update category"
// @Success 200 object response.Response{errors=string,data=string} "Update Success"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 404 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/category/{id} [PATCH]
func (c *categoryController) Update(ctx *gin.Context) {
	id := &categoryRequest.CategoryGetRequest{}
	if err := ctx.ShouldBindUri(id); err != nil {
		response := response.ErrorsResponseByCode(response.IdInvalid, "Failed to process request", response.IdInvalid, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	category, categoryErr := c.CategoryRepository.GetCategory(id.Id)
	if category.ID == 0 {
		response := response.ErrorsResponseByCode(response.RecordNotFound, "Failed to process request", response.RecordNotFound, nil)
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	if categoryErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", categoryErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	input := categoryRequest.CategoryCreateOrUpdateRequest{}
	if inputErr := ctx.ShouldBindJSON(&input); inputErr != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", inputErr.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	updateCategory, updateCategoryErr := c.CategoryService.UpdateCategory(input, id.Id)
	if updateCategoryErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", updateCategoryErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Update Success", updateCategory)
	ctx.JSON(http.StatusOK, response)
	return
}

// @Summary "Delete a single category"
// @Tags	"Category"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header	string	true	"example:Bearer token (Bearer+space+token)."	default(Bearer )
// @Param	id				path	integer	true	"Category ID"									minimum(1)
// @Success 200 object response.Response{errors=string,data=string} "Delete Success"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 404 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/category/{id} [delete]
func (c *categoryController) Delete(ctx *gin.Context) {
	input := &categoryRequest.CategoryGetRequest{}
	if err := ctx.ShouldBindUri(input); err != nil {
		response := response.ErrorsResponseByCode(response.IdInvalid, "Failed to process request", response.IdInvalid, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	category, categoryErr := c.CategoryRepository.GetCategory(input.Id)
	if category.ID == 0 {
		response := response.ErrorsResponseByCode(response.RecordNotFound, "Failed to process request", response.RecordNotFound, nil)
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	if categoryErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", categoryErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if deleteCategoryErr := c.CategoryRepository.DeleteCategory(input.Id); deleteCategoryErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", deleteCategoryErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Delete Success", nil)
	ctx.JSON(http.StatusOK, response)
	return
}
