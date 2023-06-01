package main

import (
	"fmt"
	_ "go-todolist-aws/docs"
	"go-todolist-aws/router"
	"go-todolist-aws/router/authRouter"
	"go-todolist-aws/router/categoryRouter"
	"go-todolist-aws/utils/gorm"
	"go-todolist-aws/utils/log"
	"go-todolist-aws/utils/redis"
	"net/http"
	"os"

	"github.com/apex/gateway"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go-ToDoList AWS System API
// @version 1.0
// @description AWS lambda API
// @host localhost:9753
// @BasePath /api/v1
// schemes http
func main() {
	mode := os.Getenv("GIN_MODE")
	db, err := gorm.InitMySQL()
	if err != nil {
		log.Error(err)
		return
	}

	rdb, rerr := redis.InitRedis()
	if rerr != nil {
		log.Error(rerr)
	}

	defer gorm.Close(db)
	defer redis.Close(rdb)

	r := router.Default()
	r = authRouter.GetRoute(r, db, rdb)
	r = categoryRouter.GetRoute(r, db, rdb)
	swagger := ginSwagger.URL(fmt.Sprintf("http://localhost:9753/api/swagger/doc.json"))
	r.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swagger))

	if mode == "release" {
		log.Fatal(gateway.ListenAndServe(":9753", r))
	} else {
		log.Fatal(http.ListenAndServe(":9753", r))
	}
}
