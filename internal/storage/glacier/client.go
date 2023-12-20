package glacier

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsG "github.com/aws/aws-sdk-go-v2/service/glacier"
	"github.com/wyll-io/dicomizer/internal/storage"
)

type Client struct {
	client *awsG.Client
}

func NewClient(cfg aws.Config) storage.StorageAction {
	return &Client{
		client: awsG.NewFromConfig(cfg),
	}
}
