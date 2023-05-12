package auth

import (
	"go-todolist-aws/config"
	"go-todolist-aws/model"
	"go-todolist-aws/request/authRequest"
	"go-todolist-aws/service/authService"
	"go-todolist-aws/service/jwtService"
	"go-todolist-aws/utils/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	AuthService authService.AuthService
	JwtService  jwtService.JwtService
}

func New(db *gorm.DB, rdb *redis.Client) AuthController {
	return &authController{
		AuthService: authService.New(db),
		JwtService:  jwtService.New(rdb),
	}
}

func (c *authController) Login(ctx *gin.Context) {
	input := &authRequest.LoginRequest{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	loginResult := c.AuthService.VerifyCredential(input.Email, input.Password)
	if v, ok := loginResult.(model.User); ok {
		generatedToken, generatedTokenErr := c.JwtService.GenerateToken(v.ID, time.Now().Add(time.Duration(config.JWTttl)*time.Second))
		if generatedTokenErr != nil {
			response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", generatedTokenErr.Error(), nil)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		v.Token = generatedToken
		response := response.SuccessResponse(http.StatusOK, "Login successfully", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
}

func (c *authController) Register(ctx *gin.Context) {
	input := authRequest.RegisterRequest{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	createUser, err := c.AuthService.CreateUser(input)
	if err != nil {
		if err.Error() == response.Messages[response.EmailAlreadyExists] {
			response := response.ErrorsResponseByCode(http.StatusInternalServerError, "Failed to process request", response.EmailAlreadyExists, nil)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Register successfully", createUser)
	ctx.JSON(http.StatusOK, response)
	return
}
