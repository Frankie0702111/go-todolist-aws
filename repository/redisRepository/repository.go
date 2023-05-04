package redisRepository

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	Set(key string, value interface{}, expire time.Duration) (string, error)
	Get(key string) (interface{}, error)
	Del(key string) (interface{}, error)
	GetInt(key string) (int, error)
	IncrBy(key string, value int64) (uint64, error)
	ExpireAt(key string, time time.Time) bool
}

type redisRepository struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) RedisRepository {
	return &redisRepository{
		rdb: rdb,
	}
}
