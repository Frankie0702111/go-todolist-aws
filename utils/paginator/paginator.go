package paginator

import (
	"gorm.io/gorm"
)

// Standard paging structure that receives the most primitive DO (Data Object)
// It is recommended to create another structure like this externally and convert DO (Data Object) to DTO (Data Transfer Object) or VO (Value Object)
type Page[T any] struct {
	CurrentPage int64 `json:"currentPage"`
	PageLimit   int64 `json:"pageLimit"`
	Total       int64 `json:"total"` // Data count
	Pages       int64 `json:"pages"` // Total page
	Data        []T   `json:"data"`
}

// First set up the query, then call this func
func (page *Page[T]) SelectPages(query *gorm.DB) (e error) {
	var model T
	query.Model(&model).Count(&page.Total)
	if page.Total == 0 {
		page.Data = []T{}
		return
	}

	e = query.Model(&model).Scopes(Paginate(page)).Find(&page.Data).Error
	return
}

func Paginate[T any](page *Page[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.CurrentPage < 1 {
			page.CurrentPage = 1
		}

		switch {
		// Limit the maximum number of pagination
		case page.PageLimit > 100:
			page.PageLimit = 100
		// Set the default number of pagination
		case page.PageLimit <= 0:
			page.PageLimit = 10
		}

		page.Pages = page.Total / page.PageLimit
		if page.Total%page.PageLimit != 0 {
			page.Pages++
		}

		p := page.CurrentPage
		if page.CurrentPage > page.Pages {
			p = page.Pages
		}

		size := page.PageLimit
		offset := int((p - 1) * size)
		return db.Offset(offset).Limit(int(size))
	}
}
