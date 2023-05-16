package authRepository_test

import (
	"fmt"
	"go-todolist-aws/config"
	"go-todolist-aws/model"
	"go-todolist-aws/repository/authRepository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func setUp(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s_test?charset=utf8mb4&parseTime=True&loc=Local", config.SourceUser, config.SourcePassword, config.TestSourceHost1, config.TestSourcePort1, config.SourceDataBase)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
}

func TestVerifyCredential_Success(t *testing.T) {
	setUp(t)

	user := model.User{
		Username: "Test123",
		Email:    "test123@test.com",
		Password: "password",
	}
	createUser := db.Create(&user)
	assert.NoError(t, createUser.Error)

	r := authRepository.New(db)

	res := r.VerifyCredential(user.Email)
	assert.NotNil(t, res)
	assert.Equal(t, user.Email, res.(model.User).Email)

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}

// Test for an invalid credential
func TestVerifyCredential_Failed(t *testing.T) {
	setUp(t)

	r := authRepository.New(db)
	res := r.VerifyCredential("invalid@example.com")
	assert.Nil(t, res)
}

func TestFindByEmail_Success(t *testing.T) {
	setUp(t)

	user := model.User{
		Username: "Test456",
		Email:    "test456@test.com",
		Password: "password",
	}
	createUser := db.Create(&user)
	assert.NoError(t, createUser.Error)

	r := authRepository.New(db)

	res, _ := r.FindByEmail(user.Email)
	assert.Equal(t, user.Email, res.Email)

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}

// Test for an invalid email
func TestFindByEmail_Failed(t *testing.T) {
	setUp(t)

	r := authRepository.New(db)
	res := r.VerifyCredential("invalid@example.com")
	assert.Nil(t, res)
}

func TestCreateUser_Success(t *testing.T) {
	setUp(t)

	user := model.User{
		Username: "Test789",
		Email:    "test789@test.com",
		Password: "password",
	}
	r := authRepository.New(db)

	createUser, createUserErr := r.CreateUser(user)
	assert.NoError(t, createUserErr)
	assert.NotEmpty(t, createUser.Password)
	assert.NotEqual(t, user.Password, createUser.Password)

	deleteUser := db.Where("email", user.Email).Delete(&user)
	assert.NoError(t, deleteUser.Error)
}

// Duplicate user created
func TestCreateUser_Failed(t *testing.T) {
	setUp(t)

	user := model.User{
		Username: "admin",
		Email:    "admin@test.com",
		Password: "password",
	}
	r := authRepository.New(db)
	_, createUserErr := r.CreateUser(user)
	assert.Error(t, createUserErr)
}
