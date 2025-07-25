package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func NicknameAvailabilityCheck(queryFn func() (*dynamodb.QueryOutput, error)) (bool, error) {
	// separated out so logic can be easily tested
	out, err := queryFn()

	if err != nil {
		return false, err
	}

	// IF THIS FUNCTION RETURNS TRUE THEN NICKNAME IS AVAILABLE
	isAvailable := len(out.Items) == 0

	return isAvailable, nil
}

func IsNicknameAvailableInDynamodb(nickname string, TableName string, DdbClient *dynamodb.Client, ctx context.Context) (bool, error) {
	// use this function when validating nickname from other files
	return NicknameAvailabilityCheck(func() (*dynamodb.QueryOutput, error) {
		return DdbClient.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(TableName),
			IndexName:              aws.String("NicknameIndex"),
			KeyConditionExpression: aws.String("nickname = :nick"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":nick": &types.AttributeValueMemberS{Value: nickname},
			},
			Limit: aws.Int32(1),
		})
	})
}
