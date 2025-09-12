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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserDynamoHelper struct {
	DbClient  *dynamodb.Client
	TableName string
	Ctx       context.Context
}

func (deps *UserDynamoHelper) AddUser(u *models.User) error {
	newNickname := models.NewNickname(u.Nickname, u.Name, u.Userid)
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				// adds nickname item to reserve name
				Put: &types.Put{
					TableName: aws.String(deps.TableName),
					Item:      *newNickname.DatabaseFormat(),
					// if this fails most likely because nickname is already in use everything will roll back
					ConditionExpression: aws.String("attribute_not_exists(pk)"),
				},
			},
			{
				// add user to db
				Put: &types.Put{
					TableName: aws.String(deps.TableName),
					Item:      *u.DatabaseFormat(),
				},
			},
		},
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func (deps *UserDynamoHelper) FindByNickname(nickname string) (*models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(deps.TableName),
		IndexName:              aws.String("NicknameIndex"),
		KeyConditionExpression: aws.String("nickname = :nick"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nick": &types.AttributeValueMemberS{Value: strings.ToLower(nickname)},
		},
		Limit: aws.Int32(1),
	}

	output, err := deps.DbClient.Query(deps.Ctx, input)

	if err != nil {
		return nil, err
	}

	if output.Count < 1 {
		return nil, nil
	}

	return models.ConvertToUser(output.Items[0])
}

func (deps *UserDynamoHelper) FindById(id string) (*models.User, error) {
	input := dynamodb.GetItemInput{
		Key:       models.UserKey(id),
		TableName: &deps.TableName,
	}

	output, err := deps.DbClient.GetItem(deps.Ctx, &input)

	if err != nil {
		return nil, err
	}

	if len(output.Item) < 1 {
		return nil, nil
	}

	return models.ConvertToUser(output.Item)
}

func (deps *UserDynamoHelper) DeleteFromDynamo(userId string, nickname string) error {
	// delete user profile, nickname, friends, post and allat
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				// delete account metadata
				Delete: &types.Delete{
					Key:       models.UserKey(userId),
					TableName: aws.String(deps.TableName),
				},
			},
			{
				// delete nickname reservation
				Delete: &types.Delete{
					Key:       models.NicknameKey(nickname),
					TableName: aws.String(deps.TableName),
				},
			},
		},
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			log.Println("Could not find user resource to delete: %w", nickname)
			return nil
		}
		return err
	}

	return nil
}

func (deps *UserDynamoHelper) UpdateNicknameAndFullname(userId string, nickname string, fullname string) error {

	if !utils.NameIsValid(&fullname) {
		return fmt.Errorf("Name is invalid!")
	}

	if !utils.NicknameValid(nickname) {
		return fmt.Errorf("Nickname is invalid!")
	}

	nicknameAvail, err := deps.NicknameAvailable(strings.ToLower(nickname))

	if err != nil {
		return fmt.Errorf("Error while trying to check nickname availability: "+err.Error(), nil)
	}

	if !nicknameAvail {
		return fmt.Errorf(nickname+" is already in use", nil)
	}

	// delete old nickname item
	// update user details nickname
	// put new nickname item

	newNickname := models.NewNickname(nickname, fullname, userId)

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Delete: &types.Delete{
					Key:       models.NicknameKey(nickname),
					TableName: aws.String(deps.TableName),
				},
			},
			{
				Update: &types.Update{
					Key:              models.UserKey(userId),
					TableName:        aws.String(deps.TableName),
					UpdateExpression: aws.String("SET nickname = :n"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":n": &types.AttributeValueMemberS{Value: strings.ToLower(nickname)},
					},
				},
			},
			{
				Update: &types.Update{
					Key:              models.UserKey(userId),
					TableName:        aws.String(deps.TableName),
					UpdateExpression: aws.String("SET name = :fn"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":fn": &types.AttributeValueMemberS{Value: fullname},
					},
				},
			},
			{
				Put: &types.Put{
					Item:      *newNickname.DatabaseFormat(),
					TableName: aws.String(deps.TableName),
				},
			},
		},
	}

	_, upErr := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if upErr != nil {
		return fmt.Errorf("An error occurred while trying to update nickname: "+upErr.Error(), nil)
	}

	return nil
}

func (deps *UserDynamoHelper) UpdateBio(userId string, bio string) error {
	return nil
}

func (deps *UserDynamoHelper) UpdateDpUrl(userId string, url string) error {
	return nil
}

func (deps *UserDynamoHelper) NicknameAvailable(nickname string) (bool, error) {
	input := dynamodb.GetItemInput{
		Key:       models.NicknameKey(nickname),
		TableName: aws.String(deps.TableName),
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
