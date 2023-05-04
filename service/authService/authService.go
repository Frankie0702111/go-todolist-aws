package authService

import (
	"go-todolist-aws/model"
	"go-todolist-aws/utils/log"

	"golang.org/x/crypto/bcrypt"
)

func (s *authService) VerifyCredential(email string, password string) interface{} {
	res := s.AuthRepository.VerifyCredential(email)
	if v, ok := res.(model.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}

		return false
	}

	return false
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Error("Faild to compare password : " + err.Error())
		return false
	}

	return true
}
