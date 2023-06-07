package s3Repository

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository interface {
	FileUpload(file *multipart.FileHeader, uuidV4 string) (*manager.UploadOutput, error)
	FileRemove(file string, uuidV4 string) error
}

type s3Repository struct {
	db *s3.Client
}

func New(db *s3.Client) S3Repository {
	return &s3Repository{
		db: db,
	}
}
