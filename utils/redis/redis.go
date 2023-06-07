package redis

import (
	"fmt"
	"go-todolist-aws/config"
	"go-todolist-aws/utils/log"

	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	log.Info("Testing Golang Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
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
