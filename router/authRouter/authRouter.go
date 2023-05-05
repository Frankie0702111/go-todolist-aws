package authRouter

import (
	"go-todolist-aws/controller/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoute(r *gin.Engine, db *gorm.DB) *gin.Engine {
	controller := auth.New(db)
	r.Group("api/v1/auth/login").POST("/", controller.Login)

	return r
}
