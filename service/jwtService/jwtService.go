package jwtService

import (
	"fmt"
	"go-todolist-aws/config"
	"go-todolist-aws/utils/log"
	"strconv"
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

func ValidateToken(token string) (*jwt.Token, error) {
	// Parse the token
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			// Return an error if the signing method isn't HMAC
			log.Error("Unexpected signing method : " + t_.Header["alg"].(string))
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		// Return the key
		return []byte(getSecretKey()), nil
	})
}

func (s *jwtService) GenerateToken(userID uint64, t int) (string, error) {
	now := time.Now()
	jwtTTL := time.Duration(t) * time.Second
	claims := &jwtCustomClaim{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    s.issuer,
		},
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := generateToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	_, setRedisErr := s.RedisRepository.Set("token"+strconv.FormatUint(userID, 10), token, jwtTTL)
	if setRedisErr != nil {
		return "", setRedisErr
	}

	return token, nil
}

func (s *jwtService) Logout(userID uint64) error {
	_, err := s.RedisRepository.Del("token" + strconv.FormatUint(userID, 10))
	if err != nil {
		return err
	}

	return nil
}
