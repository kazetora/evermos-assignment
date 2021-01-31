package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient redis.UniversalClient
var ctx context.Context

// Setup Setup redis client
func Setup() error {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_SERVER"),
		DB:       0,
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	ctx = context.Background()
	pong, err := redisClient.Ping(ctx).Result()

	fmt.Println(pong, err)
	return err
}

// RedisClient return redis client
func RedisClient() redis.UniversalClient {
	return redisClient
}

// Context return storage context
func Context() context.Context {
	return ctx
}
