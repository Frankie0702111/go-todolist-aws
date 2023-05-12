package authRepository

import (
	"go-todolist-aws/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	VerifyCredential(email string) interface{}
	FindByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}
