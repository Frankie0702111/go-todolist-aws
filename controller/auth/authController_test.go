package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-todolist-aws/config"
	"go-todolist-aws/model"
	"go-todolist-aws/request/authRequest"
	"go-todolist-aws/router"
	"go-todolist-aws/router/authRouter"
	"go-todolist-aws/utils/log"
	"go-todolist-aws/utils/response"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	r   *gin.Engine
	db  *gorm.DB
	rdb *redis.Client
	err error
)

func setUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s_test?charset=utf8mb4&parseTime=True&loc=Local", config.SourceUser, config.SourcePassword, config.TestSourceHost1, config.TestSourcePort1, config.SourceDataBase)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.TestSourceHost2, config.TestSourcePort2),
		Password: config.RedisPassword,
		DB:       0,
	})
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		t.Fatalf("Failed to connect redis: %v", err)
	}

	r = router.Default()
	r = authRouter.GetRoute(r, db, rdb)
}

func performRequest(method string, path string, data interface{}, token interface{}) (*httptest.ResponseRecorder, error) {
	log.Info("token", token)
	// Structure data to JSON
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create a request
	req, err := http.NewRequest(method, path, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+fmt.Sprintf("%v", token))
	}

	// Response of the capture processor
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, nil
}

func TestLogin_Success(t *testing.T) {
	setUp(t)

	// Create a test user
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

	input := authRequest.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	w, err := performRequest("POST", "/api/v1/auth/login", input, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	response := &response.Response{}
	responseJSON := json.Unmarshal(w.Body.Bytes(), response)
	assert.NoError(t, responseJSON)
	assert.Equal(t, "Login successfully", response.Message)

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}

// Invalid email and password
func TestLogin_Failed(t *testing.T) {
	setUp(t)

	user := model.User{
		Username: "Test456",
		Email:    "test456@test.com",
		Password: "",
	}
	password := "password"
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	assert.NoError(t, err)

	user.Password = string(hashPassword)
	createUser := db.Create(&user)
	assert.NoError(t, createUser.Error)

	input := authRequest.LoginRequest{
		Email:    "test789@test.com",
		Password: password,
	}

	w, err := performRequest("POST", "/api/v1/auth/login", input, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusOK, w.Code)

	input.Email = user.Email
	input.Password = "12345678"
	w, err = performRequest("POST", "/api/v1/auth/login", input, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusOK, w.Code)

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}

func TestRegister_Success(t *testing.T) {
	setUp(t)

	input := authRequest.RegisterRequest{
		Username: "Test789",
		Email:    "test789@test.com",
		Password: "password",
	}

	w, err := performRequest("POST", "/api/v1/auth/register", input, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	deleteUser := db.Where("email", input.Email).Delete(&model.User{})
	assert.NoError(t, deleteUser.Error)
}

// Duplicate user data
func TestRegister_Failed(t *testing.T) {
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

	input := authRequest.RegisterRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: password,
	}

	w, err := performRequest("POST", "/api/v1/auth/register", input, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusOK, w.Code)

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}

func TestRefreshToken_Success(t *testing.T) {
	setUp(t)

	user := model.User{
		Username: "Test456",
		Email:    "test456@test.com",
		Password: "",
	}
	password := "password"
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	assert.NoError(t, err)

	user.Password = string(hashPassword)
	createUser := db.Create(&user)
	assert.NoError(t, createUser.Error)

	input := authRequest.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	// Get token by login
	w, err := performRequest("POST", "/api/v1/auth/login", input, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	log.Info("login w.Header", w.Header())
	log.Info("login w.Body", w.Body)

	loginResponse := &response.Response{}
	loginResponseJSON := json.Unmarshal(w.Body.Bytes(), loginResponse)
	assert.NoError(t, loginResponseJSON)
	assert.Equal(t, "Login successfully", loginResponse.Message)

	loginToken := loginResponse.Data.(map[string]interface{})["token"].(string)

	// Request too fast, sleep 1 second
	time.Sleep(1 * time.Second)

	// Use the login token to refresh the token
	w, err = performRequest("POST", "/api/v1/auth/refresh", "", loginToken)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	log.Info("refresh w.Header", w.Header())
	log.Info("refresh w.Body", w.Body)

	res := &response.Response{}
	resJSON := json.Unmarshal(w.Body.Bytes(), res)
	assert.NoError(t, resJSON)
	assert.Equal(t, "Refresh token successfully", res.Message)

	refreshToken := res.Data.(map[string]interface{})["token"].(string)
	assert.NotEqual(t, loginToken, refreshToken)

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}

func TestRefreshToken_Failed(t *testing.T) {
	setUp(t)

	w, err := performRequest("POST", "/api/v1/auth/refresh", "", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJpc3MiOiJnb2p3dCIsImV4cCI6MTY4NDUwNzcwOSwiaWF0IjoxNjg0NTA2ODA5fQ.EGO1So_CpdolmnEXBTjeDeHFRao0wEUTT4vd_dVkj48")
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusOK, w.Code)

	res := &response.Response{}
	resJSON := json.Unmarshal(w.Body.Bytes(), res)
	assert.NoError(t, resJSON)
	assert.Equal(t, "Token is not valid", res.Message)
}

func TestLogout_Success(t *testing.T) {
	setUp(t)

	user := model.User{
		Username: "Test456",
		Email:    "test456@test.com",
		Password: "",
	}
	password := "password"
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	assert.NoError(t, err)

	user.Password = string(hashPassword)
	createUser := db.Create(&user)
	assert.NoError(t, createUser.Error)

	input := authRequest.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	w, err := performRequest("POST", "/api/v1/auth/login", input, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	loginResponse := &response.Response{}
	loginResponseJSON := json.Unmarshal(w.Body.Bytes(), loginResponse)
	assert.NoError(t, loginResponseJSON)
	assert.Equal(t, "Login successfully", loginResponse.Message)

	loginToken := loginResponse.Data.(map[string]interface{})["token"].(string)
	time.Sleep(1 * time.Second)

	w, err = performRequest("POST", "/api/v1/auth/logout", "", loginToken)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	res := &response.Response{}
	resJSON := json.Unmarshal(w.Body.Bytes(), res)
	assert.NoError(t, resJSON)
	assert.Equal(t, "Successfully logged out", res.Message)

	createUser = db.Delete(&user)
	assert.NoError(t, createUser.Error)
}

func TestLogout_Failed(t *testing.T) {
	setUp(t)

	w, err := performRequest("POST", "/api/v1/auth/logout", "", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJpc3MiOiJnb2p3dCIsImV4cCI6MTY4NDUwNzcwOSwiaWF0IjoxNjg0NTA2ODA5fQ.EGO1So_CpdolmnEXBTjeDeHFRao0wEUTT4vd_dVkj48")
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusOK, w.Code)

	res := &response.Response{}
	resJSON := json.Unmarshal(w.Body.Bytes(), res)
	assert.NoError(t, resJSON)
	assert.Equal(t, "Token is not valid", res.Message)
}
