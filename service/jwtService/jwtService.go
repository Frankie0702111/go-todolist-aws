package jwtService

import (
	"go-todolist-aws/config"
	"go-todolist-aws/utils/log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtCustomClaim struct {
	// The userId is the only required field
	UserID uint64 `json:"user_id"`

	// This is a registered JWT claim (StandardClaims are deprecated)
	jwt.RegisteredClaims
}

func getSecretKey() string {
	// Get the secret key from the environment variable
	secretKey := config.JWTSecretKey
	if secretKey == "" {
		// If the environment variable is empty, use a default value
		secretKey = "learnGolangJWTToken"
	}
	return secretKey
}

func (s *jwtService) GenerateToken(userID uint64, t time.Time) string {
	claims := &jwtCustomClaim{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := generateToken.SignedString([]byte(s.secretKey))
	if err != nil {
		log.Error("Failed to process request : Signature failed")
		return ""
	}

	// _, setRedisErr := s.RedisRepository.Set("token"+strconv.FormatUint(userID, 10), token, time.Duration(config.JWTttl)*time.Second)
	// if setRedisErr != nil {
	// 	log.Error("Failed to set the token in redis : " + setRedisErr.Error())
	// }

	return token
}
