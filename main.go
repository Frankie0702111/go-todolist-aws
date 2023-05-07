package main

import (
	"go-todolist-aws/router"
	"go-todolist-aws/router/authRouter"
	"go-todolist-aws/router/categoryRouter"
	"go-todolist-aws/utils/gorm"
	"go-todolist-aws/utils/log"
	"net/http"
	"os"

	"github.com/apex/gateway"
)

func main() {
	mode := os.Getenv("GIN_MODE")
	db, err := gorm.InitMySQL()
	if err != nil {
		log.Error(err)
		return
	}

	defer gorm.Close(db)

	r := router.Default()
	r = authRouter.GetRoute(r, db)
	r = categoryRouter.GetRoute(r, db)

	if mode == "release" {
		log.Fatal(gateway.ListenAndServe(":9753", r))
	} else {
		log.Fatal(http.ListenAndServe(":9753", r))
	}
}
