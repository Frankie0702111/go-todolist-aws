package jwtService_test

import (
	"fmt"
	"go-todolist-aws/config"
	"go-todolist-aws/repository/redisRepository"
	"go-todolist-aws/service/jwtService"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var (
	rdb *redis.Client
	err error
)

func setUp(t *testing.T) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.TestSourceHost2, config.TestSourcePort2),
		Password: config.RedisPassword,
		DB:       0,
	})
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		t.Fatalf("Failed to connect redis: %v", err)
	}
}

func TestGenerateToken_Success(t *testing.T) {
	setUp(t)

	s := jwtService.New(rdb)
	r := redisRepository.New(rdb)

	userID := uint64(1)
	token, err := s.GenerateToken(userID, time.Now().Add(time.Duration(config.JWTttl)*time.Second))
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	res, err := r.Get("token" + strconv.FormatUint(userID, 10))
	assert.NoError(t, err)
	assert.Equal(t, token, res)
}
