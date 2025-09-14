package helpers

import (
	"breadcrumb-backend-go/models"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SearchDynamoHelper struct {
	DbClient  *dynamodb.Client
	TableName string
	Ctx       context.Context
}

func (deps *SearchDynamoHelper) SearchUser(searchStr string) ([]models.UserSearch, error) {

	if len(searchStr[:models.UserSearchIndexPrefixLen]) < models.UserSearchIndexPrefixLen {
		return nil, fmt.Errorf("Search string is too short!")
	}

	input := dynamodb.QueryInput{
		TableName: aws.String(deps.TableName),
		KeyConditions: map[string]types.Condition{
			"pk": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: models.UserSearchIndexPkPrefix + searchStr[:models.UserSearchIndexPrefixLen]},
				},
			},
			"sk": {
				ComparisonOperator: types.ComparisonOperatorContains,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: searchStr},
				},
			},
		},
	}

	result, err := deps.DbClient.Query(deps.Ctx, &input)

	if err != nil {
		return nil, err
	}

	var users []models.UserSearch

	if err := attributevalue.UnmarshalListOfMaps(result.Items, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (deps *SearchDynamoHelper) AddUserSearchIndex(user *models.User) error {
	// Adds items to search table to allow for queries where search string is similar to nickname or full name
	builder := models.UserSearch{
		UserId:   user.Userid,
		Nickname: user.Nickname,
		Name:     user.Name,
		DpUrl:    user.DpUrl,
	}

	indexes, err := builder.BuildSearchIndexes()

	if err != nil {
		return err
	}

	// creates slice of items
	var items []types.TransactWriteItem
	for _, index := range indexes {
		items = append(items, types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(deps.TableName),
				Item:      index,
			},
		})
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	}

	_, dbErr := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if dbErr != nil {
		// Check for transaction cancellation reasons
		var tce *types.TransactionCanceledException
		if errors.As(dbErr, &tce) {
			for i, reason := range tce.CancellationReasons {
				fmt.Printf("add user search index Cancellation %d: Code=%s, Message=%s\n",
					i,
					aws.ToString(reason.Code),
					aws.ToString(reason.Message),
				)
			}
		}
		return dbErr
	}

	return nil
}

func (deps *SearchDynamoHelper) DeleteUserIndexes(user *models.User) error {
	// rebuild indexes, then query them and delete
	builder := models.UserSearch{
		UserId:   user.Userid,
		Nickname: user.Nickname,
		Name:     user.Name,
		DpUrl:    user.DpUrl,
	}

	indexes, builderErr := builder.BuildSearchIndexes()
	if builderErr != nil {
		return builderErr
	}

	keys := models.GetUserSearchIndexesKeys(indexes)
	var items []types.TransactWriteItem
	for _, key := range keys {
		items = append(items, types.TransactWriteItem{
			Delete: &types.Delete{
				TableName: aws.String(deps.TableName),
				Key:       key,
			},
		})
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	}

	log.Println(keys)

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)
	if err != nil {
		// Check for transaction cancellation reasons
		var tce *types.TransactionCanceledException
		if errors.As(err, &tce) {
			for i, reason := range tce.CancellationReasons {
				fmt.Printf("delete user index Cancellation %d: Code=%s, Message=%s\n",
					i,
					aws.ToString(reason.Code),
					aws.ToString(reason.Message),
				)
			}
		}
		return err
	}

	return nil
}
