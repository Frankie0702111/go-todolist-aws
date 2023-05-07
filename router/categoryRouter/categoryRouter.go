package categoryRouter

import (
	"go-todolist-aws/controller/category"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(r *gin.Engine, db *gorm.DB) *gin.Engine {
	controller := category.New(db)
	r.Group("api/v1/category").GET("/:id", controller.Get)

	return r
}
