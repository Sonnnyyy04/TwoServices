package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "testTwoServices/proto/sonyyy04.user.v1"
)

func NewUserServiceClient(addr string) (proto.UserServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return proto.NewUserServiceClient(conn), nil
}
