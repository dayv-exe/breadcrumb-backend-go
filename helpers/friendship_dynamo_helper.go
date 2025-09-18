package helpers

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FriendshipDynamoHelper struct {
	DbClient  *dynamodb.Client
	TableName string
	Ctx       context.Context
}

func (deps *FriendshipDynamoHelper) SendFriendReq(senderId string, recipientId string) error {
	item, marshalErr := models.NewFriendRequest(recipientId, senderId).DatabaseFormat()

	if marshalErr != nil {
		log.Println("An error occurred while trying to convert friendship item to dynamodb item")
		return marshalErr
	}

	input := &dynamodb.PutItemInput{
		Item:      *item,
		TableName: aws.String(deps.TableName),
	}

	_, putErr := deps.DbClient.PutItem(deps.Ctx, input)

	if putErr != nil {
		log.Print("An error occurred while trying to put friendship item in the db")
		return putErr
	}

	return nil
}

func (deps *FriendshipDynamoHelper) CancelFriendRequest(senderId, recipientId string) error {
	input := &dynamodb.DeleteItemInput{
		Key:       models.FriendRequestKey(recipientId, senderId),
		TableName: aws.String(deps.TableName),
	}

	_, err := deps.DbClient.DeleteItem(deps.Ctx, input)
	if err != nil {
		log.Print("error while trying to delete friend request item")
		return err
	}

	return nil
}

func (deps *FriendshipDynamoHelper) EndFriendship(user1id, user2id string) error {
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Delete: &types.Delete{
					Key:       models.FriendKey(user1id, user2id),
					TableName: aws.String(deps.TableName),
				},
			},
			{
				Delete: &types.Delete{
					Key:       models.FriendKey(user2id, user1id),
					TableName: aws.String(deps.TableName),
				},
			},
		},
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)
	if err != nil {
		log.Print("error while trying to transact write (delete) user friendships")
		// Check for transaction cancellation reasons
		var tce *types.TransactionCanceledException
		if errors.As(err, &tce) {
			for i, reason := range tce.CancellationReasons {
				fmt.Printf("add user Cancellation %d: Code=%s, Message=%s\n",
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

func (deps *FriendshipDynamoHelper) AcceptFriendRequest(senderId, recipientId string) error {
	// 2 items for easier lookups
	item1, err1 := models.NewFriendship(senderId, recipientId).DatabaseFormat()
	if err1 != nil {
		log.Print("an error occurred while marshalling new friendship")
		return err1
	}
	item2, err2 := models.NewFriendship(recipientId, senderId).DatabaseFormat()
	if err2 != nil {
		log.Print("an error occurred while marshalling new friendship")
		return err2
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				// deletes friend request
				Delete: &types.Delete{
					Key:       models.FriendRequestKey(recipientId, senderId),
					TableName: aws.String(deps.TableName),
				},
			},

			// makes them friends
			{
				Put: &types.Put{
					Item:      *item1,
					TableName: aws.String(deps.TableName),
				},
			},
			{
				Put: &types.Put{
					Item:      *item2,
					TableName: aws.String(deps.TableName),
				},
			},
		},
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if err != nil {
		log.Print("an error occurred while transact writing new friendship to db")
		// Check for transaction cancellation reasons
		utils.PrintTransactWriteCancellationReason(err)
		return err
	}

	return nil
}

func (deps *FriendshipDynamoHelper) RejectFriendRequest(senderId, recipientId string) error {
	input := &dynamodb.DeleteItemInput{
		Key:       models.FriendRequestKey(recipientId, senderId),
		TableName: aws.String(deps.TableName),
	}

	_, err := deps.DbClient.DeleteItem(deps.Ctx, input)

	if err != nil {
		log.Print("an error occurred while trying to delete friend request")
		return err
	}

	return nil
}

func (deps *FriendshipDynamoHelper) usersAreFriends(senderId string, recipientId string) (bool, error) {
	input := &dynamodb.GetItemInput{
		Key:       models.FriendKey(senderId, recipientId),
		TableName: aws.String(deps.TableName),
	}

	item, err := deps.DbClient.GetItem(deps.Ctx, input)

	if err != nil {
		log.Print("an error occurred while trying to get friendship item from db")
		return false, err
	}

	return len(item.Item) > 0, nil
}

func (deps *FriendshipDynamoHelper) userHasRequestedFriendship(senderId string, recipientId string) (bool, error) {
	input := &dynamodb.GetItemInput{
		Key:       models.FriendRequestKey(recipientId, senderId),
		TableName: aws.String(deps.TableName),
	}

	item, err := deps.DbClient.GetItem(deps.Ctx, input)

	if err != nil {
		log.Println("an error occurred while trying to get friend request item from the database")
		return false, err
	}
	return len(item.Item) > 0, nil
}

func (deps *FriendshipDynamoHelper) GetFriendshipStatus(senderId string, recipientId string) (string, error) {
	requested, reqErr := deps.userHasRequestedFriendship(senderId, recipientId)
	if reqErr != nil {
		return "", reqErr
	}

	if requested {
		return constants.FRIENDSHIP_STATUS_REQUESTED, nil
	}

	received, recErr := deps.userHasRequestedFriendship(recipientId, senderId)
	if recErr != nil {
		log.Print("error while checking if user has RECEIVED a friend request")
		return "", recErr
	}

	if received {
		return constants.FRIENDSHIP_STATUS_RECEIVED, nil
	}

	friends, fErr := deps.usersAreFriends(senderId, recipientId)
	if fErr != nil {
		return "", fErr
	}

	if friends {
		return constants.FRIENDSHIP_STATUS_FRIENDS, nil
	}

	return constants.FRIENDSHIP_STATUS_NOT_FRIENDS, nil
}

func (deps *FriendshipDynamoHelper) GetAllFriends(userId string, limit int32) (*[]models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(deps.TableName),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: utils.AddPrefix(models.FriendItemPk, userId)},
			":skPrefix": &types.AttributeValueMemberS{Value: models.FriendItemSk},
		},
		Limit: aws.Int32(limit),
	}

	items, err := deps.DbClient.Query(deps.Ctx, input)
	if err != nil {
		log.Print("an error occurred while trying to query users friends")
		return nil, err
	}

	var friends []models.User
	if mErr := attributevalue.UnmarshalListOfMaps(items.Items, &friends); mErr != nil {
		log.Print("an error occurred while trying to unmarshal list of users friends")
		return nil, mErr
	}

	return &friends, nil
}

func (deps *FriendshipDynamoHelper) GetAllFriendRequests(userId string, limit int32) (*[]models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(deps.TableName),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: utils.AddPrefix(models.FriendRequestPkPrefix, userId)},
			":skPrefix": &types.AttributeValueMemberS{Value: models.FriendRequestSkPrefix},
		},
		Limit: aws.Int32(limit),
	}

	items, err := deps.DbClient.Query(deps.Ctx, input)
	if err != nil {
		log.Print("an error occurred while trying to query users friend requests")
		return nil, err
	}

	var friends []models.User
	if mErr := attributevalue.UnmarshalListOfMaps(items.Items, &friends); mErr != nil {
		log.Print("an error occurred while trying to unmarshal list of users friend requests")
		return nil, mErr
	}

	return &friends, nil
}
