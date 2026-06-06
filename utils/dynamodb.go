package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("us-east-2"),
	)

	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}
