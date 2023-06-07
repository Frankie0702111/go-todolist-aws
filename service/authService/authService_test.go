package authService_test

import (
	"fmt"
	"go-todolist-aws/config"
	"go-todolist-aws/model"
	"go-todolist-aws/request/authRequest"
	"go-todolist-aws/service/authService"
	"go-todolist-aws/utils/response"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func setUp(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s_test?charset=utf8mb4&parseTime=True&loc=Local", config.SourceUser, config.SourcePassword, config.SourceHost, config.SourcePort, config.SourceDataBase)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
}

func TestVerifyCredential_Success(t *testing.T) {
	t.Helper()
	setUp(t)

	user := model.User{
		Username: "Test123",
		Email:    "test123@test.com",
		Password: "",
	}
	password := "password"
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	assert.NoError(t, err)

	user.Password = string(hashPassword)
	createUser := db.Create(&user)
	assert.NoError(t, createUser.Error)

	s := authService.New(db)
	res := s.VerifyCredential(user.Email, password)
	assert.IsType(t, model.User{}, res)
	assert.Equal(t, user.ID, res.(model.User).ID)
	assert.Equal(t, user.Email, res.(model.User).Email)
	assert.Equal(t, user.Password, res.(model.User).Password)

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}

// Invalid email and password
func TestVerifyCredential_Failed(t *testing.T) {
	setUp(t)

	s := authService.New(db)
	res := s.VerifyCredential("test@example.com", "password")

	// Use if to assert that Error and Equal have the same result
	if res != false {
		assert.Error(t, res.(error))
	}
	assert.Equal(t, false, res)

	res = s.VerifyCredential("admin@test.com", "invalid-password")
	assert.Equal(t, false, res)
}

func TestCreateUser_Success(t *testing.T) {
	setUp(t)

	input := authRequest.RegisterRequest{
		Username: "Test456",
		Email:    "test456@test.com",
		Password: "password",
	}

	s := authService.New(db)
	res, err := s.CreateUser(input)
	assert.NoError(t, err)
	assert.NotZero(t, res.ID)
	assert.Equal(t, input.Username, res.Username)
	assert.Equal(t, input.Email, res.Email)

	deleteUser := db.Where("email", input.Email).Delete(&res)
	assert.NoError(t, deleteUser.Error)
}

// Email already exists
func TestCreateUser_Failed(t *testing.T) {
	setUp(t)

	user := model.User{
		Username: "Test789",
		Email:    "test789@test.com",
		Password: "password",
	}
	createUser := db.Create(&user)
	assert.NoError(t, createUser.Error)

	s := authService.New(db)
	_, err := s.CreateUser(authRequest.RegisterRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	assert.Equal(t, response.Messages[response.EmailAlreadyExists], err.Error())

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}
