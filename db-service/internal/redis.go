package internal

import "github.com/redis/go-redis/v9"

func NewRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	return rdb
}
