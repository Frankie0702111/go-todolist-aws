package authRepository

import (
	"go-todolist-aws/model"

	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (r *authRepository) VerifyCredential(email string) interface{} {
	var user model.User
	if res := r.db.Where("email = ?", email).Take(&user); res.Error == nil {
		return user
	}

	return nil
}

func (r *authRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	if res := r.db.Where("email = ?", email).Take(&user); res.Error != nil {
		return user, res.Error
	}

	return user, nil
}

func (r *authRepository) CreateUser(user model.User) (model.User, error) {
	var err error
	user.Password, err = hashAndSalt([]byte(user.Password))
	if err != nil {
		return user, err
	}

	if res := r.db.Save(&user); res.Error != nil {
		return user, res.Error
	}

	return user, nil
}
