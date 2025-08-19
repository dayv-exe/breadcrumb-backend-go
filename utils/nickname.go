package utils

import (
	"breadcrumb-backend-go/constants"
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func NicknameValid(nickname string) bool {
	if len(nickname) < constants.MIN_USERNAME_CHARS || len(nickname) > constants.MAX_USERNAME_CHARS {
		return false
	}

	// rules: must start and end with letter or number, only one of dot or underscore, must be 3 to 15 characters long
	match, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9._]*[a-zA-Z0-9]$`, nickname)
	if !match {
		return false
	}

	if strings.Contains(nickname, "..") || strings.Contains(nickname, "__") {
		return false
	}

	if strings.Contains(nickname, "_.") || strings.Contains(nickname, "._") {
		return false
	}

	return true
}

func NickNameAvailable(nickname string, TableName string, DdbClient *dynamodb.Client, ctx context.Context) (bool, error) {
	// returns true if the nickname parsed is not taken in the db
	return nicknameAvailableQueryRunner(func() (*dynamodb.QueryOutput, error) {
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

func nicknameAvailableQueryRunner(queryFn func() (*dynamodb.QueryOutput, error)) (bool, error) {
	result, err := queryFn()

	if err != nil {
		return false, err
	}

	isAvailable := len(result.Items) < 1
	return isAvailable, nil
}
