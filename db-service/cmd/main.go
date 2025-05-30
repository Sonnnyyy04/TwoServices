package main

import (
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"testTwoServices/db-service/internal"
	"testTwoServices/db-service/internal/kafkaext"
	"testTwoServices/db-service/server"
	proto "testTwoServices/proto/sonyyy04.user.v1"
)

func main() {
	//init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	//init db
	pg, err := internal.NewPostgres()
	if err != nil {
		sugar.Fatalf("failed to connect to postgres", "error", err)
	}
	defer pg.Close()

	//init Redis
	redis := internal.NewRedis()
	defer redis.Close()

	//init kafka
	if err := kafkaext.EnsureKafkaTopic("kafka:9092", "events", 1); err != nil {
		sugar.Fatalw("failed to ensure kafka topic", "error", err)
	}

	// Kafka writer
	kWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "events",
	})

	//init Grpc
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		sugar.Fatalf("failed to listen", "error", err)
	}
	grpcServer := grpc.NewServer()
	userService := server.NewUserService(pg, redis, kWriter, sugar)
	proto.RegisterUserServiceServer(grpcServer, userService)

	sugar.Infow("starting gRPC server", "port", 50051)
	if err := grpcServer.Serve(listener); err != nil {
		sugar.Fatalw("failed to serve gRPC", "error", err)
	}
}
