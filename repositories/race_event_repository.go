package repositories

import (
	"context"

	"f1-race-events/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type RaceEventRepository struct {
	client    *dynamodb.Client
	tableName string
}

func InitRaceEventRepository(
	client *dynamodb.Client,
	tableName string,
) *RaceEventRepository {

	return &RaceEventRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *RaceEventRepository) Create(
	ctx context.Context,
	event models.RaceEvent,
) error {

	item, err := attributevalue.MarshalMap(event)

	if err != nil {
		return err
	}

	_, err = r.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: &r.tableName,
			Item:      item,
		},
	)

	return err
}
