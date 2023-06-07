package taskRequest

import (
	"go-todolist-aws/request"
	"mime/multipart"
	"time"
)

// const TimeFormat = "2006-01-02 15:04:05"

// type CivilTime time.Time

// func (t *CivilTime) UnmarshalJSON(data []byte) (err error) {
// 	// Null values are not parsed
// 	if len(data) == 2 {
// 		*t = CivilTime(time.Time{})
// 		return nil
// 	}

// 	// Specify the format of parsing
// 	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
// 	if err != nil {
// 		return err
// 	}

// 	*t = CivilTime(now)
// 	return nil
// }

// // Handling c.JSON parsed value issues
// func (t CivilTime) MarshalJSON() ([]byte, error) {
// 	b := make([]byte, 0, len(TimeFormat)+2)
// 	b = append(b, '"')
// 	b = time.Time(t).AppendFormat(b, TimeFormat)
// 	b = append(b, '"')
// 	return b, nil
// }

// // Called when writing to MySQL
// func (t CivilTime) Value() (driver.Value, error) {
// 	// 0001-01-01 00:00:00 is a null value, and resolves directly to a null value when it is encountered
// 	if t.String() == "0001-01-01 00:00:00" {
// 		return nil, nil
// 	}
// 	return []byte(time.Time(t).Format(TimeFormat)), nil
// }

// // Called when searching MySQL
// func (t *CivilTime) Scan(v interface{}) error {
// 	// The format of mysql internal date is 2006-01-02 15:04:05 +0000 UTC, so it needs to be formatted again when searching
// 	tTime, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", v.(time.Time).String())
// 	*t = CivilTime(tTime)
// 	return nil
// }

// // For "fmt.Println" and subsequent validation scenarios
// func (t CivilTime) String() string {
// 	return time.Time(t).Format(TimeFormat)
// }

type TaskGetListRequest struct {
	Id              int64      `form:"id" json:"id,omitempty"`
	UserID          int64      `form:"user_id" json:"user_id,omitempty"`
	Title           string     `form:"title" json:"title,omitempty" binding:"max=100"`
	SpecifyDatetime *time.Time `form:"specify_datetime" json:"specify_datetime,omitempty" time_format:"2006-01-02 15:04:05"`
	IsSpecifyTime   *bool      `form:"is_specify_time" json:"is_specify_time,omitempty"`
	IsComplete      *bool      `form:"is_complete" json:"is_complete,omitempty"`
	request.Pagination
}

type TaskCreateRequest struct {
	UserID          int64                 `form:"user_id" json:"user_id" binding:"required"`
	CategoryID      int64                 `form:"category_id" json:"category_id" binding:"required"`
	Title           string                `form:"title" json:"title" binding:"required,max=100"`
	Note            string                `form:"note" json:"note,omitempty"`
	Url             string                `form:"url" json:"url,omitempty"`
	Image           *multipart.FileHeader `form:"image" json:"image,omitempty"`
	SpecifyDatetime *time.Time            `form:"specify_datetime" json:"specify_datetime,omitempty" time_format:"2006-01-02 15:04:05"`
	IsSpecifyTime   bool                  `form:"is_specify_time" json:"is_specify_time,omitempty"`
	Priority        int8                  `form:"priority" json:"priority" binding:"required,oneof=1 2 3"`
	IsComplete      bool                  `form:"is_complete" json:"is_complete,omitempty"`
}

type TaskUpdateRequest struct {
	CategoryID      int64                 `form:"category_id" json:"category_id,omitempty"`
	Title           string                `form:"title" json:"title,omitempty" binding:"max=100"`
	Note            string                `form:"note" json:"note,omitempty"`
	Url             string                `form:"url" json:"url,omitempty"`
	Image           *multipart.FileHeader `form:"image" json:"image,omitempty"`
	SpecifyDatetime *time.Time            `form:"specify_datetime" json:"specify_datetime,omitempty" time_format:"2006-01-02 15:04:05"`
	IsSpecifyTime   bool                  `form:"is_specify_time" json:"is_specify_time,omitempty"`
	Priority        int8                  `form:"priority" json:"priority" binding:"required,oneof=1 2 3"`
	IsComplete      bool                  `form:"is_complete" json:"is_complete,omitempty"`
}

type TaskGetRequest struct {
	request.TableID
}
