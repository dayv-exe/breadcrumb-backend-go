package utils

// helper to put stuff in db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PutItemInDDbTable(newItem map[string]types.AttributeValue, tableName string, ddbClient *dynamodb.Client, ctx context.Context) (*dynamodb.PutItemOutput, error) {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      newItem,
	}

	out, err := ddbClient.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to write user to DynamoDB: %w", err)
	}

	return out, nil
}
