package auth

import (
	"go-todolist-aws/config"
	"go-todolist-aws/model"
	"go-todolist-aws/request/authRequest"
	"go-todolist-aws/service/authService"
	"go-todolist-aws/service/jwtService"
	"go-todolist-aws/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	Logout(ctx *gin.Context)
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

// Login is a function for user login
// @Summary "User Login"
// @Tags	"Auth"
// @Version 1.0
// @Produce application/json
// @Param	* body authRequest.LoginRequest true "User Login"
// @Success 200 object response.Response{errors=string,data=string} "Login successfully"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 401 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/auth/login [post]
func (c *authController) Login(ctx *gin.Context) {
	input := &authRequest.LoginRequest{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		response := response.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	loginResult := c.AuthService.VerifyCredential(input.Email, input.Password)
	if v, ok := loginResult.(model.User); ok {
		generatedToken, generatedTokenErr := c.JwtService.GenerateToken(v.ID, config.JWTttl)
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

	response := response.ErrorsResponseByCode(http.StatusUnauthorized, "Failed to process request", response.InvalidCredential, nil)
	ctx.JSON(http.StatusUnauthorized, response)
	return
}

// Register is a function for user register
// @Summary "User Register"
// @Tags	"Auth"
// @Version 1.0
// @Produce application/json
// @Param	* body authRequest.RegisterRequest true "User Register"
// @Success 201 object response.Response{errors=string,data=string} "Register Success"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "Failed to process request"
// @Router	/auth/register [post]
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

// RefreshToken is a function for token refresh
// @Summary "User Refresh Token"
// @Tags	"Auth"
// @Version 1.0
// @Produce application/json
// @Param	Authorization header string true "example:Bearer token (Bearer+space+token)." default(Bearer )
// @Success 200 object response.Response{errors=string,data=string} "Login successfully"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 401 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "(Token is invalid in the server. Token is not valid. Failed to process request)"
// @Router	/auth/refresh [post]
func (c *authController) RefreshToken(ctx *gin.Context) {
	tokenInfo, exists := ctx.Get("token_info")
	if !exists {
		response := response.ErrorsResponseByCode(http.StatusInternalServerError, "Token is not valid", response.TokenBindingHasUnknownErrors, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userTokenInfo, ok := tokenInfo.(*model.TokenInfo)
	if !ok {
		response := response.ErrorsResponseByCode(http.StatusInternalServerError, "Token is not valid", response.DataBindingHasUnknownErrors, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	token, err := c.JwtService.GenerateToken(userTokenInfo.UserID, config.JWTttl)
	if err != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Refresh token successfully", map[string]string{"token": token})
	ctx.JSON(http.StatusOK, response)
	return
}

// Logout is a function for user logout
// @Summary "User Logout"
// @Tags	"Auth"
// @Version 1.0
// @Produce application/json
// @Param	Authorization header string true "example:Bearer token (Bearer+space+token)." default(Bearer )
// @Success 200 object response.Response{errors=string,data=string} "Successfully logged out"
// @Failure 400 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 401 object response.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object response.Response{errors=string,data=string} "(Token is invalid in the server. Token is not valid. Failed to process request)"
// @Router	/auth/logout [post]
func (c *authController) Logout(ctx *gin.Context) {
	tokenInfo, exists := ctx.Get("token_info")
	if !exists {
		response := response.ErrorsResponseByCode(http.StatusInternalServerError, "Token is not valid", response.TokenBindingHasUnknownErrors, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userTokenInfo, ok := tokenInfo.(*model.TokenInfo)
	if !ok {
		response := response.ErrorsResponseByCode(http.StatusInternalServerError, "Token is not valid", response.DataBindingHasUnknownErrors, nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	err := c.JwtService.Logout(userTokenInfo.UserID)
	if err != nil {
		response := response.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.SuccessResponse(http.StatusOK, "Successfully logged out", nil)
	ctx.JSON(http.StatusOK, response)
	return
}
