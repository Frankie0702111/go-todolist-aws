package authRepository

import (
	"go-todolist-aws/model"
)

func (r *authRepository) VerifyCredential(email string) interface{} {
	var user model.User
	res := r.db.Where("email = ?", email).Take(&user)

	if res.Error == nil {
		return user
	}

	return nil
}
