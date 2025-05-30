package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"testTwoServices/api-service/client"
	"testTwoServices/api-service/handler"
)

func main() {
	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	// init grpc client
	grpcClient, err := client.NewUserServiceClient("db:50051")
	if err != nil {
		sugar.Fatalw("failed to connect to db-service", "error", err)
	}

	//init gin
	r := gin.Default()
	h := handler.NewHandler(grpcClient, sugar)
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	// start service
	logger.Info("Starting API server on :8080")
	if err = r.Run(":8080"); err != nil {
		sugar.Fatalw("failed run API server", "error", err)
	}
}
