package s3Repository

import (
	"context"
	"errors"
	"go-todolist-aws/config"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (r *s3Repository) FileUpload(file *multipart.FileHeader, uuidV4 string) (*manager.UploadOutput, error) {
	client := r.db
	if client == nil {
		return nil, errors.New("Invalid credential.")
	}

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		// 10 MiB
		u.PartSize = 10 * 1024 * 1024
	})

	f, openErr := file.Open()
	if openErr != nil {
		return nil, openErr
	}

	result, resulterr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(config.AWS_BUCKET),
		Key:    aws.String(uuidV4 + "/" + file.Filename),
		Body:   f,
		// ACL:    "public-read",
	})

	if resulterr != nil {
		return nil, resulterr
	}

	return result, nil
}

func (r *s3Repository) FileRemove(file string, uuidV4 string) error {
	var objectIds []types.ObjectIdentifier
	client := r.db
	if client == nil {
		return errors.New("Invalid credential.")
	}

	objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(uuidV4 + "/" + file)})
	_, err := client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(config.AWS_BUCKET),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		return err
	}

	return nil
}
