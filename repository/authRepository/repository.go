package authRepository

import "gorm.io/gorm"

type AuthRepository interface {
	VerifyCredential(email string) interface{}
}

type authRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}
