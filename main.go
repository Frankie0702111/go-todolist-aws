package main

import (
	// "context"
	"go-todolist-aws/router"
	"go-todolist-aws/router/authRouter"
	"go-todolist-aws/utils/gorm"
	"os"

	"go-todolist-aws/utils/log"
	"net/http"

	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	// ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/apex/gateway"
)

// var ginLambda *ginadapter.GinLambda

func main() {
	mode := os.Getenv("GIN_MODE")
	db, err := gorm.InitMySQL()
	if err != nil {
		log.Error(err)
		return
	}

	// defer gorm.Close(db)

	r := router.Default()
	r = authRouter.GetRoute(r, db)

	if mode == "release" {
		log.Fatal(gateway.ListenAndServe(":9753", r))

		// ginLambda = ginadapter.New(r)
		// lambda.Start(Handler)
	} else {
		log.Fatal(http.ListenAndServe(":9753", r))
	}
}

// func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	return ginLambda.ProxyWithContext(ctx, request)
// }
