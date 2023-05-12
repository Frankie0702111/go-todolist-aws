package authRouter

import (
	"go-todolist-aws/controller/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func GetRoute(r *gin.Engine, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	controller := auth.New(db, rdb)
	v1 := r.Group("api/v1/auth")
	{
		v1.POST("/login", controller.Login)
		v1.POST("/register", controller.Register)
	}

	return r
}
