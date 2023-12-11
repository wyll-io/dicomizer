package glacier

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsG "github.com/aws/aws-sdk-go-v2/service/glacier"
)

type Client struct {
	client *awsG.Client
}

func NewClient(cfg aws.Config) *Client {
	return &Client{
		client: awsG.NewFromConfig(cfg),
	}
}
