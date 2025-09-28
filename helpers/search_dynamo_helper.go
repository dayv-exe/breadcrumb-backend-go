package helpers

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SearchDynamoHelper struct {
	DbClient        *dynamodb.Client
	Ctx             context.Context
	SearchTableName string
}

func (deps *SearchDynamoHelper) SearchUser(searchStr string, limit int32) (*[]models.UserDisplayInfo, error) {

	var matches []models.UserDisplayInfo
	seen := make(map[string]int)

	tokens := utils.SplitOnDelimiter(strings.ToLower(utils.NormalizeString(searchStr)), " ", "_", ".") // splits the search string into tokens

	for _, token := range tokens {
		if len(token) >= models.UserSearchIndexPrefixLen {
			input := dynamodb.QueryInput{
				TableName:              aws.String(deps.SearchTableName),
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

			usersFound, ufErr := models.SearchItemsToUserInfoStruct(&found.Items)
			if ufErr != nil {
				log.Print("error while converting search items to user")
				return nil, ufErr
			}

			matches = append(matches, (*usersFound)...)
		}
	}

	// loop through matches and rank them and put them in results
	var result []models.UserDisplayInfo
	for index, user := range matches {
		user.Userid = strings.TrimPrefix(user.Userid, models.UserPkPrefix) // removes the 'USER#'
		key := user.Userid
		_, ok := seen[key]
		if !ok {
			// first time seen
			seen[key] = index
			result = append(result, user)
		}
	}

	return &result, nil
}

func (deps *SearchDynamoHelper) GetUserSearchIndexItems(user *models.User) ([]types.TransactWriteItem, error) {
	// Adds items to search table to allow for queries where search string is similar to nickname or full name
	builder := models.UserSearch{
		UserDisplayInfo: models.UserDisplayInfo{
			Userid:                  user.Userid,
			Nickname:                user.Nickname,
			Name:                    user.Name,
			DpUrl:                   user.DpUrl,
			DefaultProfilePicColors: user.DefaultProfilePicColors,
		},
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
				TableName: aws.String(deps.SearchTableName),
				Item:      index,
			},
		})
	}

	return items, nil
}

func (deps *SearchDynamoHelper) GetDeleteUserIndexesItems(user *models.User) ([]types.TransactWriteItem, error) {
	// rebuild indexes, then query them and get their primary keys
	builder := models.UserSearch{
		UserDisplayInfo: models.UserDisplayInfo{
			Userid:                  user.Userid,
			Nickname:                user.Nickname,
			Name:                    user.Name,
			DpUrl:                   user.DpUrl,
			DefaultProfilePicColors: user.DefaultProfilePicColors,
		},
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
				TableName: aws.String(deps.SearchTableName),
				Key:       key,
			},
		})
	}

	log.Println(items)

	return items, nil
}
