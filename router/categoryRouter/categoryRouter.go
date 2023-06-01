package categoryRouter

import (
	"go-todolist-aws/controller/category"
	"go-todolist-aws/middleware"
	"go-todolist-aws/repository/redisRepository"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func GetRoute(r *gin.Engine, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	controller := category.New(db)
	redisRepository := redisRepository.New(rdb)

	auth := r.Group("api/v1", middleware.Verify(redisRepository))
	{
		auth.POST("/category", controller.Create)
		auth.GET("/category", controller.GetByList)
		auth.GET("/category/:id", controller.Get)
	}

	return r
}
