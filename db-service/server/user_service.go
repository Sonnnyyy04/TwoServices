package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"strings"
	proto "testTwoServices/proto/sonyyy04.user.v1"
)

type UserService struct {
	db     *sql.DB
	redis  *redis.Client
	kafka  *kafka.Writer
	logger *zap.SugaredLogger
	proto.UnimplementedUserServiceServer
}

func NewUserService(db *sql.DB, r *redis.Client, k *kafka.Writer, l *zap.SugaredLogger) *UserService {
	return &UserService{
		db:     db,
		redis:  r,
		kafka:  k,
		logger: l,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	s.logger.Infow("CreateUser called", "name", req.Name, "email", req.Email)
	var id string
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users(name, email) VALUES ($1, $2) RETURNING id",
		req.Name, req.Email).Scan(&id)
	if err != nil {
		s.logger.Errorw("failed to insert into db", "error", err)
		return nil, err
	}

	//redis
	key := fmt.Sprintf("user:%s", id)
	value := fmt.Sprintf("%s|%s", req.Name, req.Email)
	err = s.redis.Set(ctx, key, value, 0).Err()
	if err != nil {
		s.logger.Errorw("failed insert into redis", "error", err)
		return nil, err
	}

	//kafka
	err = s.kafka.WriteMessages(ctx, kafka.Message{
		Key:   []byte(id),
		Value: []byte("UserCreated"),
	})
	if err != nil {
		s.logger.Errorw("failed send in kafka", "err", err)
		return nil, err
	}

	return &proto.CreateUserResponse{Id: id}, err
}

func (s *UserService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	s.logger.Infow("GetUser called", "id", req.Id)
	key := fmt.Sprintf("user:%s", req.Id)
	val, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		s.logger.Infow("value hit", "value", val)
		var name, email string
		parts := strings.SplitN(val, "|", 2)
		if len(parts) == 2 {
			name = parts[0]
			email = parts[1]
		}
		return &proto.GetUserResponse{
			Id:    req.Id,
			Name:  name,
			Email: email,
		}, nil
	} else if err != redis.Nil {
		s.logger.Errorw("redis error", "error", err)
	}
	s.logger.Infow("value miss, querying db")
	var name, email string
	err = s.db.QueryRowContext(ctx, "SELECT name, email FROM users WHERE id = $1", req.Id).Scan(&name, &email)
	if err != nil {
		s.logger.Errorw("failed to get user from db", "err", err)
		return nil, err
	}
	val = fmt.Sprintf("%s|%s", name, email)
	_ = s.redis.Set(ctx, key, val, 0).Err()
	return &proto.GetUserResponse{
		Id:    req.Id,
		Name:  name,
		Email: email,
	}, nil
}
