package authService

import (
	"errors"
	"go-todolist-aws/model"
	"go-todolist-aws/request/authRequest"
	"go-todolist-aws/utils/log"
	"go-todolist-aws/utils/response"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Error("Faild to compare password : " + err.Error())
		return false
	}

	return true
}

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

func (s *authService) CreateUser(user authRequest.RegisterRequest) (model.User, error) {
	createUser := model.User{}
	if err := smapping.FillStruct(&createUser, smapping.MapFields(&user)); err != nil {
		return createUser, err
	}

	_, err := s.AuthRepository.FindByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res, err := s.AuthRepository.CreateUser(createUser)
			if err != nil {
				return res, err
			}

			return res, nil
		}
	}

	return createUser, errors.New(response.Messages[response.EmailAlreadyExists])
}
