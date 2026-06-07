package repositories

import (
	"context"
	"strconv"
	"strings"
	"time"

	"f1-race-events/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

type RaceEventQueryFilters struct {
	RaceID        string
	DriverID      string
	ScuderiaID    string
	EventType     string
	Lap           int
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}

func (r *RaceEventRepository) GetByRaceID(
	ctx context.Context,
	filters RaceEventQueryFilters,
) ([]models.RaceEvent, error) {
	exprNames := map[string]string{
		"#raceId": "raceId",
	}
	exprValues := map[string]types.AttributeValue{
		":raceId": &types.AttributeValueMemberS{Value: filters.RaceID},
	}
	filterExpr := make([]string, 0, 6)

	if filters.DriverID != "" {
		exprNames["#driverId"] = "driverId"
		exprValues[":driverId"] = &types.AttributeValueMemberS{Value: filters.DriverID}
		filterExpr = append(filterExpr, "#driverId = :driverId")
	}

	if filters.ScuderiaID != "" {
		exprNames["#scuderiaId"] = "scuderiaId"
		exprValues[":scuderiaId"] = &types.AttributeValueMemberS{Value: filters.ScuderiaID}
		filterExpr = append(filterExpr, "#scuderiaId = :scuderiaId")
	}

	if filters.EventType != "" {
		exprNames["#eventType"] = "eventType"
		exprValues[":eventType"] = &types.AttributeValueMemberS{Value: filters.EventType}
		filterExpr = append(filterExpr, "#eventType = :eventType")
	}

	if filters.Lap != 0 {
		exprNames["#lap"] = "lap"
		exprValues[":lap"] = &types.AttributeValueMemberN{Value: strconv.Itoa(filters.Lap)}
		filterExpr = append(filterExpr, "#lap = :lap")
	}

	if filters.CreatedAtFrom != nil && filters.CreatedAtTo != nil {
		exprNames["#createdAt"] = "createdAt"
		exprValues[":createdAtFrom"] = &types.AttributeValueMemberS{Value: filters.CreatedAtFrom.UTC().Format(time.RFC3339Nano)}
		exprValues[":createdAtTo"] = &types.AttributeValueMemberS{Value: filters.CreatedAtTo.UTC().Format(time.RFC3339Nano)}
		filterExpr = append(filterExpr, "#createdAt BETWEEN :createdAtFrom AND :createdAtTo")
	} else if filters.CreatedAtFrom != nil {
		exprNames["#createdAt"] = "createdAt"
		exprValues[":createdAtFrom"] = &types.AttributeValueMemberS{Value: filters.CreatedAtFrom.UTC().Format(time.RFC3339Nano)}
		filterExpr = append(filterExpr, "#createdAt >= :createdAtFrom")
	} else if filters.CreatedAtTo != nil {
		exprNames["#createdAt"] = "createdAt"
		exprValues[":createdAtTo"] = &types.AttributeValueMemberS{Value: filters.CreatedAtTo.UTC().Format(time.RFC3339Nano)}
		filterExpr = append(filterExpr, "#createdAt <= :createdAtTo")
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 &r.tableName,
		KeyConditionExpression:    aws.String("#raceId = :raceId"),
		ExpressionAttributeNames:  exprNames,
		ExpressionAttributeValues: exprValues,
	}

	if len(filterExpr) > 0 {
		queryInput.FilterExpression = aws.String(strings.Join(filterExpr, " AND "))
	}

	result, err := r.client.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	var events []models.RaceEvent
	if err = attributevalue.UnmarshalListOfMaps(result.Items, &events); err != nil {
		return nil, err
	}

	return events, nil
}
