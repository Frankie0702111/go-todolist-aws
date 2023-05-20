package model

import (
	"time"
)

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Token     string    `gorm:"-" json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokenInfo struct {
	UserID  uint64 `json:"user_id"`
	IssUser string `json:"iss_user"`
	Token   string `json:"token"`
	Valid   bool   `json:"valid"`
}
