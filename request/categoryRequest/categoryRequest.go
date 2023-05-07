package categoryReqyest

type CategoryGetRequest struct {
	Id int64 `uri:"id" binding:"required"`
}
