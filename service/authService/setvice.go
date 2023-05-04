package authService

import (
	"go-todolist-aws/Repository/authRepository"

	"gorm.io/gorm"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
}

type authService struct {
	AuthRepository authRepository.AuthRepository
}

func New(db *gorm.DB) AuthService {
	return &authService{
		AuthRepository: authRepository.New(db),
	}
}
