package taskRouter

import (
	"go-todolist-aws/controller/task"
	"go-todolist-aws/middleware"
	"go-todolist-aws/repository/redisRepository"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func GetRoute(r *gin.Engine, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	controller := task.New(db)
	redisRepository := redisRepository.New(rdb)

	auth := r.Group("api/v1", middleware.Verify(redisRepository))
	{
		auth.POST("/task", controller.Create)
		auth.GET("/task", controller.GetByList)
		auth.GET("/task/:id", controller.Get)
		auth.PATCH("/task/:id", controller.Update)
		auth.DELETE("/task/:id", controller.Delete)
	}

	return r
}
