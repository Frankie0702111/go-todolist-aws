package jwtService

import (
	"time"
)

type JwtService interface {
	GenerateToken(userID uint64, t time.Time) string
}

type jwtService struct {
	// Secret key used to sign the token
	secretKey string
	// Who creates the token
	issuer string
}

func New() JwtService {
	return &jwtService{
		// Call the getSecretKey function to get the secret key
		secretKey: getSecretKey(),
		// who creates the token
		issuer: "gojwt",
	}
}
