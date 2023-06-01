package request

type Pagination struct {
	// 頁數(請從1開始帶入)
	Page int64 `form:"page" json:"page" binding:"required,gt=0"`
	// 筆數(請從1開始帶入)
	Limit int64 `form:"limit" json:"limit" binding:"required,gt=0"`
}

type TableID struct {
	Id int64 `uri:"id" binding:"required"`
}
