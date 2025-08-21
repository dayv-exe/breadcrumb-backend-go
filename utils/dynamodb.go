package utils

// helper to put stuff in db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDbHelper struct {
	TableName string
	Client    *dynamodb.Client
	Ctx       *context.Context
}

func (deps DynamoDbHelper) PutItem(newItem map[string]types.AttributeValue) (*dynamodb.PutItemOutput, error) {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(deps.TableName),
		Item:      newItem,
	}

	out, err := deps.Client.PutItem(*deps.Ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to write user to DynamoDB: %w", err)
	}

	return out, nil
}
