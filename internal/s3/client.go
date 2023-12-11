package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client *awsS3.Client
}

func NewClient(cfg aws.Config) *Client {
	return &Client{
		client: awsS3.NewFromConfig(cfg),
	}
}
