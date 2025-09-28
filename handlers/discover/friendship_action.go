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

func handleSendRequest(status string, sender *models.User, recipientId string, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	if status != constants.FRIENDSHIP_STATUS_NOT_FRIENDS {
		return models.InvalidRequestErrorResponse("You have to not be friends or have any pending friend requests with this person to do this."), nil
	}
	reqErr := friendshipHelper.SendFriendReq(sender, recipientId)
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

func handleAcceptFriendship(status string, otherUser, thisUser *models.User, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	if status != constants.FRIENDSHIP_STATUS_RECEIVED {
		return models.InvalidRequestErrorResponse("You cant accept a friend request from someone who never sent you one"), nil
	}
	reqErr := friendshipHelper.AcceptFriendRequest(otherUser, thisUser)
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

func handleGetAllFriends(thisUserId string, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	allFriends, afErr := friendshipHelper.GetAllFriends(thisUserId)
	if afErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while trying to get all friends, try again.", afErr, "error in friendship action handler"), nil
	}

	return models.SuccessfulGetRequestResponse(allFriends), nil
}

func handleGetAllFriendRequests(thisUserId string, friendshipHelper helpers.FriendshipDynamoHelper) (events.APIGatewayProxyResponse, error) {
	allReqs, afErr := friendshipHelper.GetAllFriendRequests(thisUserId)
	if afErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while trying to get all friend requests, try again.", afErr, "error in friendship action handler"), nil
	}

	return models.SuccessfulGetRequestResponse(allReqs), nil
}

func (deps *FriendRequestDependencies) HandleFriendshipAction(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// /friendship/{action}/{userid}
	action, actionExists := req.PathParameters["action"]
	if !actionExists {
		return models.InvalidRequestErrorResponse(""), nil
	}

	otherUserId, exists := req.PathParameters["userid"]
	if !exists && (action != constants.FRIENDSHIP_ACTION_GET_FRIENDS && action != constants.FRIENDSHIP_ACTION_GET_REQUESTED) {
		return models.InvalidRequestErrorResponse(""), nil
	}

	thisUserId := utils.GetAuthUserId(req)
	if thisUserId == "" {
		return models.UnauthorizedErrorResponse(""), nil
	}

	// check if they are friends
	friendshipHelper := helpers.FriendshipDynamoHelper{
		DbClient:  deps.DbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	status, statusErr := friendshipHelper.GetFriendshipStatus(thisUserId, otherUserId)
	if statusErr != nil {
		return models.ServerSideErrorResponse("Something went wrong, try again.", statusErr, "error while trying to get friendship status"), nil
	}

	findUserHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	thisUser, tuErr := findUserHelper.FindById(thisUserId)
	if tuErr != nil {
		return models.ServerSideErrorResponse("Something went wrong.", tuErr, "error while trying to fetch this user details"), nil
	}

	otherUser, ouErr := findUserHelper.FindById(otherUserId)
	if ouErr != nil {
		return models.ServerSideErrorResponse("Something went wrong", tuErr, "error while trying to fetch other user details"), nil
	}

	switch action {
	// send friend request if not friends
	case constants.FRIENDSHIP_ACTION_REQUEST:
		return handleSendRequest(status, thisUser, otherUserId, friendshipHelper)

		// cancel friend request if sent
	case constants.FRIENDSHIP_ACTION_CANCEL_REQUEST:
		return handleCancelRequest(status, thisUserId, otherUserId, friendshipHelper)

		// end friendship if friends
	case constants.FRIENDSHIP_ACTION_END_FRIENDSHIP:
		return handleEndFriendship(status, thisUserId, otherUserId, friendshipHelper)

	// accept friend request
	case constants.FRIENDSHIP_ACTION_ACCEPT:
		return handleAcceptFriendship(status, otherUser, thisUser, friendshipHelper)

		// reject friend request
	case constants.FRIENDSHIP_ACTION_REJECT:
		return handleRejectFriendship(status, otherUserId, thisUserId, friendshipHelper)

		// to list all friends
	case constants.FRIENDSHIP_ACTION_GET_FRIENDS:
		return handleGetAllFriends(thisUserId, friendshipHelper)

		// to list all friend requests
	case constants.FRIENDSHIP_ACTION_GET_REQUESTED:
		return handleGetAllFriendRequests(thisUserId, friendshipHelper)
	default:
		return models.ServerSideErrorResponse("Something went wrong while determining friendship action, try again.", fmt.Errorf("Status returned does not match any expected outcome. status returned: %v", status), "Status returned does not match any expected outcome"), nil
	}
}
