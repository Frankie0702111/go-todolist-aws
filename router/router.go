package router

import (
	"go-todolist-aws/middleware"

	"github.com/gin-gonic/gin"
)

func Default() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	router := gin.New()
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	test := router.Group("api/v1/test")
	{
		test.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Test success.",
			})
		})
	}

	return router
}
