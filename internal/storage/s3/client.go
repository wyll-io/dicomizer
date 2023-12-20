package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/wyll-io/dicomizer/internal/storage"
)

type Client struct {
	client *awsS3.Client
}

func NewClient(cfg aws.Config) storage.StorageAction {
	return &Client{
		client: awsS3.NewFromConfig(cfg),
	}
}
