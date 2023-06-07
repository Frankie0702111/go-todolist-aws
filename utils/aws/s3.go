package aws

import (
	"context"
	localConfig "go-todolist-aws/config"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(localConfig.AWS_ACCESS_KEY_ID, localConfig.AWS_SECRET_ACCESS_KEY, "")),
		config.WithRegion(localConfig.AWS_REGION),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
