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
	"gorm.io/gorm"
)

type AuthController interface {
	Login(ctx *gin.Context)
}

type authcontroller struct {
	AuthService authService.AuthService
	JwtService  jwtService.JwtService
}

func New(db *gorm.DB) AuthController {
	return &authcontroller{
		AuthService: authService.New(db),
		JwtService:  jwtService.New(),
	}
}

func (c *authcontroller) Login(ctx *gin.Context) {
	input := &authRequest.LoginRequest{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	loginResult := c.AuthService.VerifyCredential(input.Email, input.Password)
	if v, ok := loginResult.(model.User); ok {
		generatedToken := c.JwtService.GenerateToken(v.ID, time.Now().Add(time.Duration(config.JWTttl)*time.Second))
		if len(generatedToken) < 1 {
			response := response.ErrorsResponseByCode(http.StatusInternalServerError, "Failed to process request", response.SignatureFailed, nil)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		v.Token = generatedToken
		response := response.SuccessResponse(http.StatusOK, "Login successfully", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
}
