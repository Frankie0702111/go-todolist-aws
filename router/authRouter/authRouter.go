package authRouter

import (
	"go-todolist-aws/controller/auth"
	"go-todolist-aws/middleware"
	"go-todolist-aws/repository/redisRepository"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func GetRoute(r *gin.Engine, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	controller := auth.New(db, rdb)
	redisRepository := redisRepository.New(rdb)
	public := r.Group("api/v1/auth")
	{
		public.POST("/login", controller.Login)
		public.POST("/register", controller.Register)
	}

	auth := r.Group("api/v1/auth", middleware.Verify(redisRepository))
	{
		auth.POST("/refresh", controller.RefreshToken)
		auth.POST("/logout", controller.Logout)
	}

	return r
}
