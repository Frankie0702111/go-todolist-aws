package model

import (
	"time"
)

type Task struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	CategoryID      int64      `json:"category_id"`
	Category        Category   `gorm:"foreignkey:CategoryID;references:ID" json:"category"`
	Title           string     `json:"title"`
	Note            string     `json:"note"`
	Url             string     `json:"url"`
	Img             string     `json:"img"`
	ImgLink         string     `json:"img_link"`
	ImgUuid         string     `json:"img_uuid"`
	SpecifyDatetime *time.Time `json:"specify_datetime"`
	IsSpecifyTime   bool       `json:"is_specify_time"`
	Priority        int8       `json:"priority"`
	IsComplete      bool       `json:"is_complete"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
