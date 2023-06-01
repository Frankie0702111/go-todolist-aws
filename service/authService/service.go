package authService

import (
	"go-todolist-aws/Repository/authRepository"
	"go-todolist-aws/model"
	"go-todolist-aws/request/authRequest"

	"gorm.io/gorm"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user authRequest.RegisterRequest) (model.User, error)
}

type authService struct {
	AuthRepository authRepository.AuthRepository
}

func New(db *gorm.DB) AuthService {
	return &authService{
		AuthRepository: authRepository.New(db),
	}
}
