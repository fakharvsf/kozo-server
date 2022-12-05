package utils

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

var RedisCtx = context.Background()
var RedisClient *redis.Client

func RedisConnect() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := RedisClient.Ping(RedisCtx).Result()

	// DB Connection
	if err != nil {
		fmt.Println("Could not connect to Redis Client:", err)
		return err
	} else {
		fmt.Println("Connected to Redis Client.")
		return nil
	}
}


