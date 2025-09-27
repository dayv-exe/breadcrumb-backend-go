package helpers

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FriendshipDynamoHelper struct {
	DbClient   *dynamodb.Client
	TableNames *utils.TableNames
	Ctx        context.Context
}

func (deps *FriendshipDynamoHelper) SendFriendReq(sender *models.User, recipientId string) error {
	item, marshalErr := models.NewFriendRequest(recipientId, sender).DatabaseFormat()

	if marshalErr != nil {
		log.Println("An error occurred while trying to convert friendship item to dynamodb item")
		return marshalErr
	}

	input := &dynamodb.PutItemInput{
		Item:      *item,
		TableName: aws.String(deps.TableNames.Users),
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
		TableName: aws.String(deps.TableNames.Users),
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
					TableName: aws.String(deps.TableNames.Users),
				},
			},
			{
				Delete: &types.Delete{
					Key:       models.FriendKey(user2id, user1id),
					TableName: aws.String(deps.TableNames.Users),
				},
			},
		},
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)
	if err != nil {
		log.Print("error while trying to transact write (delete) user friendships")
		// Check for transaction cancellation reasons
		utils.PrintTransactWriteCancellationReason(err)
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
					TableName: aws.String(deps.TableNames.Users),
				},
			},

			// makes them friends
			{
				Put: &types.Put{
					Item:      *item1,
					TableName: aws.String(deps.TableNames.Users),
				},
			},
			{
				Put: &types.Put{
					Item:      *item2,
					TableName: aws.String(deps.TableNames.Users),
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
		TableName: aws.String(deps.TableNames.Users),
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
		TableName: aws.String(deps.TableNames.Users),
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
		TableName: aws.String(deps.TableNames.Users),
	}

	item, err := deps.DbClient.GetItem(deps.Ctx, input)

	if err != nil {
		log.Println("an error occurred while trying to get friend request item from the database")
		return false, err
	}
	return len(item.Item) > 0, nil
}

func (deps *FriendshipDynamoHelper) GetFriendshipStatus(senderId string, recipientId string) (string, error) {
	// checks if this user has sent a friend request to other user
	requested, reqErr := deps.userHasRequestedFriendship(senderId, recipientId)
	if reqErr != nil {
		return "", reqErr
	}

	if requested {
		return constants.FRIENDSHIP_STATUS_REQUESTED, nil
	}

	// checks if other user has sent a friend request to this user
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
		TableName:              aws.String(deps.TableNames.Users),
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

	friends, fErr := models.FriendFormat(&items.Items)
	if fErr != nil {
		return nil, fErr
	}

	return friends, nil
}

func (deps *FriendshipDynamoHelper) GetAllFriendRequests(userId string, limit int32) (*[]models.FriendRequest, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(deps.TableNames.Users),
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

	var requests []models.FriendRequest
	if mErr := attributevalue.UnmarshalListOfMaps(items.Items, &requests); mErr != nil {
		log.Print("an error occurred while trying to unmarshal list of users friend requests")
		return nil, mErr
	}

	return &requests, nil
}
