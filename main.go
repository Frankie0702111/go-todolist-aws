package main

import (
	"context"
	"go-todolist-aws/config"
	"go-todolist-aws/router"
	"go-todolist-aws/router/authRouter"
	"go-todolist-aws/utils/gorm"
	"go-todolist-aws/utils/log"
	"go-todolist-aws/utils/redis"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	// "github.com/apex/gateway"
)

var ginLambda *ginadapter.GinLambda

func main() {
	db, err := gorm.InitMySQL()
	if err != nil {
		log.Panic(err)
	}

	rdb, rerr := redis.InitRedis()
	if rerr != nil {
		log.Panic(rerr)
	}

	defer gorm.Close(db)
	defer redis.Close(rdb)

	r := router.Default()
	r = authRouter.GetRoute(r, db, rdb)

	if config.Environtment == "production" {
		// log.Fatal(gateway.ListenAndServe(":9753", r))

		ginLambda = ginadapter.New(r)
		lambda.Start(Handler)
	} else {
		log.Fatal(http.ListenAndServe(":9753", r))
	}
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}
