package router

import "github.com/gin-gonic/gin"

func Default() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS

	return router
}
