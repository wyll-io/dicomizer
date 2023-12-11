package database

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type DB struct {
	Client *rds.Client
}

func New(cfg aws.Config) DB {
	return DB{
		Client: rds.NewFromConfig(cfg),
	}
}
