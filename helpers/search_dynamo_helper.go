package helpers

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

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

func (deps *SearchDynamoHelper) SearchUser(searchStr string, limit int32) ([]models.UserSearch, error) {

	var matches []models.UserSearch
	seen := make(map[string]int)

	tokens := utils.SplitOnDelimiter(strings.ToLower(utils.NormalizeString(searchStr)), " ", "_", ".") // splits the search string into tokens

	for _, token := range tokens {
		if len(token) >= models.UserSearchIndexPrefixLen {
			input := dynamodb.QueryInput{
				TableName:              aws.String(deps.TableName),
				KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":pk":       &types.AttributeValueMemberS{Value: models.UserSearchIndexPkPrefix + token[:models.UserSearchIndexPrefixLen]},
					":skPrefix": &types.AttributeValueMemberS{Value: token},
				},
				Limit: aws.Int32(limit),
			}

			found, qErr := deps.DbClient.Query(deps.Ctx, &input)
			if qErr != nil {
				log.Println("An error occurred inside loop for querying tokens gotten from search string")
				return nil, qErr
			}

			var usersFound []models.UserSearch
			if marshalErr := attributevalue.UnmarshalListOfMaps(found.Items, &usersFound); marshalErr != nil {
				return nil, marshalErr
			}

			matches = append(matches, usersFound...)
		}
	}

	// loop through matches and rank them and put them in results
	var result []models.UserSearch
	for index, user := range matches {
		user.UserId = strings.TrimPrefix(user.UserId, models.UserPkPrefix) // removes the 'USER#'
		key := user.UserId
		ogIndex, ok := seen[key]
		if !ok {
			// first time seen
			seen[key] = index
			result = append(result, user)
		} else {
			// seen before, then remove it adn add 1 to the rating where we first saw it
			result[ogIndex].Rating += 1
		}
	}

	return result, nil
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
