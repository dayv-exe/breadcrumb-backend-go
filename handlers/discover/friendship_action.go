package discover

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type FriendRequestDependencies struct {
	DbClient  *dynamodb.Client
	TableName string
}

func handleSendRequest(status, senderId, recipientId string, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	if status != constants.FRIENDSHIP_STATUS_NOT_FRIENDS {
		return models.InvalidRequestErrorResponse("You have to not be friends or have any pending friend requests with this person to do this."), nil
	}
	reqErr := friendshipHelper.SendFriendReq(senderId, recipientId)
	if reqErr != nil {
		return models.ServerSideErrorResponse("Something went wrong, try again.", reqErr, "Error in friend req handler"), nil
	}

	return models.SuccessfulRequestResponse("Request sent", false), nil
}

func handleCancelRequest(status, senderId, recipientId string, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	if status != constants.FRIENDSHIP_STATUS_REQUESTED {
		return models.InvalidRequestErrorResponse("You dont have any pending friend request with this person to cancel"), nil
	}
	reqErr := friendshipHelper.CancelFriendRequest(senderId, recipientId)
	if reqErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while canceling friend request, try again.", reqErr, "Error in friend req handler"), nil
	}

	return models.SuccessfulRequestResponse("Request canceled", false), nil
}

func handleEndFriendship(status, senderId, recipientId string, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	if status != constants.FRIENDSHIP_STATUS_FRIENDS {
		return models.InvalidRequestErrorResponse("You are not friends with this person to begin with, so why are you trying to end the friendship??? now your account will have to be suspended for suspicious behavior :("), nil
	}
	reqErr := friendshipHelper.EndFriendship(senderId, recipientId)
	if reqErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while ending friendship, try again.", reqErr, "Error in friend req handler"), nil
	}

	return models.SuccessfulRequestResponse("Friendship ended", false), nil
}

func handleAcceptFriendship(status, senderId, recipientId string, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	if status != constants.FRIENDSHIP_STATUS_RECEIVED {
		return models.InvalidRequestErrorResponse("You cant accept a friend request from someone who never sent you one"), nil
	}
	reqErr := friendshipHelper.AcceptFriendRequest(senderId, recipientId)
	if reqErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while accepting friendship, try again.", reqErr, "Error in friend req handler"), nil
	}

	return models.SuccessfulRequestResponse("Friendship started", false), nil
}

func handleRejectFriendship(status, senderId, recipientId string, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	if status != constants.FRIENDSHIP_STATUS_RECEIVED {
		return models.InvalidRequestErrorResponse("You cant reject a friend request from someone who never sent you one"), nil
	}
	reqErr := friendshipHelper.RejectFriendRequest(senderId, recipientId)
	if reqErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while rejecting friendship, try again.", reqErr, "Error in friend req handler"), nil
	}

	return models.SuccessfulRequestResponse("Friendship ended before it began", false), nil
}

func (deps *FriendRequestDependencies) HandleFriendshipAction(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	action, actionExists := req.PathParameters["action"]
	if !actionExists {
		return models.InvalidRequestErrorResponse(""), nil
	}

	recipientId, exists := req.PathParameters["userid"]
	if !exists {
		return models.InvalidRequestErrorResponse(""), nil
	}

	senderId := utils.GetAuthUserId(req)
	if senderId == "" {
		return models.UnauthorizedErrorResponse(""), nil
	}

	// check if they are friends
	friendshipHelper := helpers.FriendshipDynamoHelper{
		DbClient:  deps.DbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	status, statusErr := friendshipHelper.GetFriendshipStatus(senderId, recipientId)
	if statusErr != nil {
		return models.ServerSideErrorResponse("Something went wrong, try again.", statusErr, "error while trying to get friendship status"), nil
	}

	switch action {
	// send friend request if not friends
	case constants.FRIENDSHIP_ACTION_REQUEST:
		return handleSendRequest(status, senderId, recipientId, friendshipHelper)

		// cancel friend request if sent
	case constants.FRIENDSHIP_ACTION_CANCEL_REQUEST:
		return handleCancelRequest(status, senderId, recipientId, friendshipHelper)

		// end friendship if friends
	case constants.FRIENDSHIP_ACTION_END_FRIENDSHIP:
		return handleEndFriendship(status, senderId, recipientId, friendshipHelper)

	// accept friend request
	case constants.FRIENDSHIP_ACTION_ACCEPT:
		return handleAcceptFriendship(status, senderId, recipientId, friendshipHelper)

		// reject friend request
	case constants.FRIENDSHIP_ACTION_REJECT:
		return handleRejectFriendship(status, senderId, recipientId, friendshipHelper)

	default:
		return models.ServerSideErrorResponse("Something went wrong while determining friendship action, try again.", fmt.Errorf("Status returned does not match any expected outcome. status returned: %v", status), "Status returned does not match any expected outcome"), nil
	}
}
