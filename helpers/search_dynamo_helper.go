package helpers

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SearchDynamoHelper struct {
	DbClient   *dynamodb.Client
	Ctx        context.Context
	TableNames *utils.TableNames
}

func (deps *SearchDynamoHelper) SearchUser(searchStr string, limit int32) ([]models.UserSearch, error) {

	var matches []models.UserSearch
	seen := make(map[string]int)

	tokens := utils.SplitOnDelimiter(strings.ToLower(utils.NormalizeString(searchStr)), " ", "_", ".") // splits the search string into tokens

	for _, token := range tokens {
		if len(token) >= models.UserSearchIndexPrefixLen {
			input := dynamodb.QueryInput{
				TableName:              aws.String(deps.TableNames.Search),
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

func (deps *SearchDynamoHelper) GetUserSearchIndexItems(user *models.User) ([]types.TransactWriteItem, error) {
	// Adds items to search table to allow for queries where search string is similar to nickname or full name
	builder := models.UserSearch{
		UserId:   user.Userid,
		Nickname: user.Nickname,
		Name:     user.Name,
		DpUrl:    user.DpUrl,
	}

	indexes, err := builder.BuildSearchIndexes()

	if err != nil {
		return nil, err
	}

	// creates slice of items
	var items []types.TransactWriteItem
	for _, index := range indexes {
		items = append(items, types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(deps.TableNames.Search),
				Item:      index,
			},
		})
	}
	return items, nil
}

func (deps *SearchDynamoHelper) GetDeleteUserIndexesItems(user *models.User) ([]types.TransactWriteItem, error) {
	// rebuild indexes, then query them and get their primary keys
	builder := models.UserSearch{
		UserId:   user.Userid,
		Nickname: user.Nickname,
		Name:     user.Name,
		DpUrl:    user.DpUrl,
	}

	indexes, builderErr := builder.BuildSearchIndexes()
	if builderErr != nil {
		return nil, builderErr
	}

	keys := models.GetUserSearchIndexesKeys(indexes)
	var items []types.TransactWriteItem
	for _, key := range keys {
		items = append(items, types.TransactWriteItem{
			Delete: &types.Delete{
				TableName: aws.String(deps.TableNames.Search),
				Key:       key,
			},
		})
	}

	log.Println(items)

	return items, nil
}
