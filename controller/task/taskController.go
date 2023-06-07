package task

import (
	"go-todolist-aws/repository/taskRepository"
	"go-todolist-aws/request/taskRequest"
	"go-todolist-aws/service/taskService"
	"go-todolist-aws/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController interface {
	Create(ctx *gin.Context)
	GetByList(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type taskController struct {
	TaskService    taskService.TaskService
	TaskRepository taskRepository.TaskRepository
}

func New(db *gorm.DB) TaskController {
	return &taskController{
		TaskService:    taskService.New(db),
		TaskRepository: taskRepository.New(db),
	}
}

// @Summary "Create task"
// @Tags	"Task"
// @Version 1.0
// @Accept	multipart/form-data
// @Produce application/json
// @Param	Authorization		header		string	true	"example:Bearer token (Bearer+space+token)."		default(Bearer )
// @Param	user_id				formData	integer	true	"User ID"											minimum(1)
// @Param	category_id			formData	integer	true	"Category ID"										minimum(1)
// @Param	title				formData	string	true	"Title"												maxLength(100)
// @Param	note				formData	string	false	"Note"
// @Param	url					formData	string	false	"Url"
// @Param	image				formData	file	false	"Image"
// @Param	specify_datetime	formData	string	false	"Specify Datetime (DateTime: 2006-01-02 15:04:05)"
// @Param	is_specify_time		formData	boolean	false	"Is Specify Time"
// @Param	priority			formData	integer	true	"Priority"											Enums(1, 2, 3) default(1)
// @Param	is_complete			formData	boolean	false	"Is Complete"										default(false)
// @Success 201 object response.Response{errors=string,data=string} "Create Success"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/task [post]
func (c *taskController) Create(ctx *gin.Context) {
	input := taskRequest.TaskCreateRequest{}
	if err := ctx.ShouldBind(&input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if input.Image != nil {
		if len(input.Image.Filename) > 100 {
			response := response.ErrorsResponseByCode(response.ImageFileNameLimitOf100, "Failed to process request", response.ImageFileNameLimitOf100, nil)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
		if input.Image.Size > (5 << 20) {
			response := response.ErrorsResponseByCode(response.ImageFileSizeLimitOf5MB, "Failed to process request", response.ImageFileSizeLimitOf5MB, nil)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
	}

	createTask, createTaskErr := c.TaskService.CreateTask(input)
	if createTaskErr != nil {
		if createTaskErr.Error() == response.Messages[response.DuplicateCreatedData] {
			response := response.ErrorsResponseByCode(response.DuplicateCreatedData, "Failed to process request", response.DuplicateCreatedData, nil)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", createTaskErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Create Success", createTask)
	ctx.JSON(http.StatusOK, response)
	return
}

// @Summary "Task list"
// @Tags	"Task"
// @Version 1.0
// @Produce application/json
// @Param	Authorization		header		string	true	"example:Bearer token (Bearer+space+token)."		default(Bearer )
// @Param	id					formData	integer	false	"Task ID"											minimum(1)
// @Param	user_id				formData	integer	false	"User ID"											minimum(1)
// @Param	title				formData	string	false	"Title"												maxLength(100)
// @Param	specify_datetime	formData	string	false	"Specify Datetime (DateTime: 2006-01-02 15:04:05)"
// @Param	is_specify_time		formData	boolean	false	"Is Specify Time"
// @Param	is_complete			formData	boolean	false	"Is Complete"
// @Param	page				query		integer	true	"Page"												minimum(1) default(1)
// @Param	limit				query		integer	true	"Limit"												minimum(2) default(5)
// @Success 200 object response.PageResponse{errors=string,data=string} "Successfully get task list"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/task [get]
func (c *taskController) GetByList(ctx *gin.Context) {
	input := taskRequest.TaskGetListRequest{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	task := c.TaskRepository.GetTaskList(input)
	response := response.SuccessPageResponse(http.StatusOK, "Successfully get task list", task.CurrentPage, task.PageLimit, task.Total, task.Pages, task.Data)
	ctx.JSON(http.StatusOK, response)
	return
}

// @Summary "Get a single task"
// @Tags	"Task"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header	string	true	"example:Bearer token (Bearer+space+token)."	default(Bearer )
// @Param	id				path	integer	true	"Task ID"										minimum(1)
// @Success 200 object response.Response{errors=string,data=string} "Record not found || Successfully get task"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/task/{id} [get]
func (c *taskController) Get(ctx *gin.Context) {
	input := &taskRequest.TaskGetRequest{}
	if err := ctx.ShouldBindUri(input); err != nil {
		response := response.ErrorsResponseByCode(response.IdInvalid, "Failed to process request", response.IdInvalid, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	task, taskErr := c.TaskRepository.GetTask(input.Id)
	if task.ID == 0 {
		response := response.SuccessResponse(http.StatusOK, "Record not found", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}
	if taskErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Successfully get task", task)
	ctx.JSON(http.StatusOK, response)
	return
}

// @Summary "Update a single task"
// @Tags	"Task"
// @Version 1.0
// @Accept	multipart/form-data
// @Produce application/json
// @Param	Authorization		header		string	true	"example:Bearer token (Bearer+space+token)."		default(Bearer )
// @Param	id					path		integer	true	"Task ID"											minimum(1)
// @Param	category_id			formData	integer	false	"Category ID"										minimum(1)
// @Param	title				formData	string	false	"Title"												maxLength(100)
// @Param	note				formData	string	false	"Note"
// @Param	url					formData	string	false	"Url"
// @Param	image				formData	file	false	"Image"
// @Param	specify_datetime	formData	string	false	"Specify Datetime (DateTime: 2006-01-02 15:04:05)"
// @Param	is_specify_time		formData	boolean	false	"Is Specify Time"
// @Param	priority			formData	integer	true	"Priority"											Enums(1, 2, 3)
// @Param	is_complete			formData	boolean	false	"Is Complete"
// @Success 200 object response.Response{errors=string,data=string} "Update Success"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 404 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/task/{id} [PATCH]
func (c *taskController) Update(ctx *gin.Context) {
	id := &taskRequest.TaskGetRequest{}
	if err := ctx.ShouldBindUri(id); err != nil {
		response := response.ErrorsResponseByCode(response.IdInvalid, "Failed to process request", response.IdInvalid, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	task, taskErr := c.TaskRepository.GetTask(id.Id)
	if task.ID == 0 {
		response := response.ErrorsResponseByCode(response.RecordNotFound, "Failed to process request", response.RecordNotFound, nil)
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	if taskErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	input := taskRequest.TaskUpdateRequest{}
	if inputErr := ctx.ShouldBind(&input); inputErr != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", inputErr.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	updateTask, updateTaskErr := c.TaskService.UpdateTask(input, task)
	if updateTaskErr != nil {
		if updateTaskErr.Error() == response.Messages[response.DuplicatedTitle] {
			response := response.ErrorsResponseByCode(response.DuplicatedTitle, "Failed to process request", response.DuplicatedTitle, nil)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", updateTaskErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Update Success", updateTask)
	ctx.JSON(http.StatusOK, response)
	return
}

// @Summary "Delete a single task"
// @Tags	"Task"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header	string	true	"example:Bearer token (Bearer+space+token)."	default(Bearer )
// @Param	id				path	integer	true	"Task ID"										minimum(1)
// @Success 200 object response.Response{errors=string,data=string} "Delete Success"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 404 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/task/{id} [delete]
func (c *taskController) Delete(ctx *gin.Context) {
	id := &taskRequest.TaskGetRequest{}
	if err := ctx.ShouldBindUri(id); err != nil {
		response := response.ErrorsResponseByCode(response.IdInvalid, "Failed to process request", response.IdInvalid, nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	task, taskErr := c.TaskRepository.GetTask(id.Id)
	if task.ID == 0 {
		response := response.ErrorsResponseByCode(response.RecordNotFound, "Failed to process request", response.RecordNotFound, nil)
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	if taskErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	deleteErr := c.TaskService.DeleteTask(task)
	if deleteErr != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", deleteErr.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Delete Success", nil)
	ctx.JSON(http.StatusOK, response)
	return
}
