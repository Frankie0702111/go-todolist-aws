package redisRepository


import (
	"context"
	"time"
)

var ctx = context.Background()

func (r *redisRepository) Set(key string, value interface{}, expire time.Duration) (string, error) {
	val, err := r.rdb.Set(ctx, key, value, expire).Result()
	return val, err
}

func (r *redisRepository) Get(key string) (interface{}, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	return val, err
}

func (r *redisRepository) Del(key string) (interface{}, error) {
	val, err := r.rdb.Del(ctx, key).Result()
	return val, err
}

func (r *redisRepository) GetInt(key string) (int, error) {
	val, err := r.rdb.Get(ctx, key).Int()
	return val, err
}

// Increase the number of requests
func (r *redisRepository) IncrBy(key string, value int64) (uint64, error) {
	val, err := r.rdb.IncrBy(ctx, key, value).Uint64()
	return val, err
}

// Expire time
func (r *redisRepository) ExpireAt(key string, time time.Time) bool {
	val := r.rdb.ExpireAt(ctx, key, time).Val()
	return val
}