package jwtService_test

import (
	"fmt"
	"go-todolist-aws/config"
	"go-todolist-aws/repository/redisRepository"
	"go-todolist-aws/service/jwtService"
	"strconv"
	"testing"

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
	token, err := s.GenerateToken(userID, config.JWTttl)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	res, err := r.Get("token" + strconv.FormatUint(userID, 10))
	assert.NoError(t, err)
	assert.Equal(t, token, res)
}

func TestLogout_Success(t *testing.T) {
	setUp(t)

	s := jwtService.New(rdb)
	r := redisRepository.New(rdb)
	userID := uint64(1)
	generateToken, err := s.GenerateToken(userID, config.JWTttl)
	assert.NoError(t, err)
	assert.NotEmpty(t, generateToken)

	deleteTokenErr := s.Logout(userID)
	assert.NoError(t, deleteTokenErr)

	res, err := r.Get("token" + strconv.FormatUint(userID, 10))
	assert.Error(t, err)
	assert.Equal(t, res, "")
}

func TestLogout_Failed(t *testing.T) {
	setUp(t)

	s := jwtService.New(rdb)
	r := redisRepository.New(rdb)
	userID := uint64(1)
	generateToken, err := s.GenerateToken(userID, config.JWTttl)
	assert.NoError(t, err)
	assert.NotEmpty(t, generateToken)

	// Delete non-existing-key
	deleteTokenErr := s.Logout(123)
	// Confirm error as nil
	assert.NoError(t, deleteTokenErr)

	// The user information still exists and the user is not properly logged out
	res, err := r.Get("token" + strconv.FormatUint(userID, 10))
	assert.NoError(t, err)
	assert.Equal(t, generateToken, res)
}
