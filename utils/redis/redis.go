package redis

import (
	"fmt"
	"go-todolist-aws/config"
	"log"

	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	log.Println("Testing Golang Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", config.RedisHost),
		Password: config.RedisPassword, // no password set
		DB:       0,                    // use default DB
	})

	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func Close(rdb *redis.Client) {
	rdb.Close()
}