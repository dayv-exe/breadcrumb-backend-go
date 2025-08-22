package utils

import (
	"breadcrumb-backend-go/constants"
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type NicknameDependencies struct {
	TableName string
	DbClient  *dynamodb.Client
	Ctx       context.Context
}

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

func (deps *NicknameDependencies) NicknameAvailable(nickname string) (bool, error) {
	input := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: strings.ToLower(nickname)},
			"sk": &types.AttributeValueMemberS{Value: "NICKNAME"},
		},
	}

	return nicknameAvailableQueryRunner(func() (*dynamodb.GetItemOutput, error) {
		return deps.DbClient.GetItem(deps.Ctx, &input)
	})
}

func nicknameAvailableQueryRunner(queryFn func() (*dynamodb.GetItemOutput, error)) (bool, error) {
	result, err := queryFn()

	if err != nil {
		return false, err
	}

	isAvailable := result.Item == nil
	return isAvailable, nil
}
