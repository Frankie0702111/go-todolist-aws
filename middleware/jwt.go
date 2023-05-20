package middleware

import (
	"go-todolist-aws/model"
	"go-todolist-aws/repository/redisRepository"
	"go-todolist-aws/service/jwtService"
	"go-todolist-aws/utils/log"
	"go-todolist-aws/utils/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Verify(r redisRepository.RedisRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response := response.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", response.TokenInvalid, nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			response := response.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", response.BearerTokenNotInProperFormat, nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		authHeader = strings.TrimSpace(splitToken[1])
		// Validate the token
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			response := response.ErrorsResponse(http.StatusUnauthorized, "Token is not valid", err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get the claims of the token
		claims := token.Claims.(jwt.MapClaims)
		// output the user_id
		log.Info("Claim[user_id]: ", claims["user_id"])
		// output the issuer
		log.Info("Claim[user_id]: ", claims["iss"])

		tokenInfo := &model.TokenInfo{
			UserID:  uint64(claims["user_id"].(float64)),
			IssUser: claims["iss"].(string),
			Token:   strings.TrimSpace(splitToken[1]),
			Valid:   token.Valid,
		}

		// Get whitelist token via redis
		redisToken, err := r.Get("token" + strconv.FormatUint(tokenInfo.UserID, 10))
		if err != nil {
			response := response.ErrorsResponse(http.StatusInternalServerError, "Token is invalid in the server", err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		// Check if the token is the same as redis (for user logout)
		if tokenInfo.Token != redisToken {
			response := response.ErrorsResponseByCode(http.StatusBadRequest, "Token is invalid in the server", response.TokenInvalid, nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		ctx.Set("token_info", tokenInfo)
		ctx.Next()
	}
}
