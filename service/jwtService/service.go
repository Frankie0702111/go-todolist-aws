package jwtService

import (
	"go-todolist-aws/repository/redisRepository"
	"time"

	"github.com/go-redis/redis/v8"
)

type JwtService interface {
	GenerateToken(userID uint64, t time.Time) (string, error)
}

type jwtService struct {
	RedisRepository redisRepository.RedisRepository
	// Secret key used to sign the token
	secretKey string
	// Who creates the token
	issuer string
}

func New(rdb *redis.Client) JwtService {
	return &jwtService{
		RedisRepository: redisRepository.New(rdb),
		// Call the getSecretKey function to get the secret key
		secretKey: getSecretKey(),
		// who creates the token
		issuer: "gojwt",
	}
}
