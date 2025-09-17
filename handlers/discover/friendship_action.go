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

func (deps *FriendRequestDependencies) HandleFriendshipAction(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// url: /friendship/act/{userid}
	// url: /friendship/accept/{userid}
	// url: /friendship/decline/{userid}
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

	switch status {
	// send friend request if not friends
	case constants.FRIENDSHIP_STATUS_NOT_FRIENDS:
		reqErr := friendshipHelper.SendFriendReq(senderId, recipientId)
		if reqErr != nil {
			return models.ServerSideErrorResponse("Something went wrong, try again.", reqErr, "Error in friend req handler"), nil
		}

		return models.SuccessfulRequestResponse("Request sent", false), nil

		// cancel friend request if sent
	case constants.FRIENDSHIP_STATUS_REQUESTED:
		reqErr := friendshipHelper.CancelFriendRequest(senderId, recipientId)
		if reqErr != nil {
			return models.ServerSideErrorResponse("Something went wrong while canceling friend request, try again.", reqErr, "Error in friend req handler"), nil
		}

		return models.SuccessfulRequestResponse("Request canceled", false), nil

		// end friendship if friends
	case constants.FRIENDSHIP_STATUS_FRIENDS:
		reqErr := friendshipHelper.EndFriendship(senderId, recipientId)
		if reqErr != nil {
			return models.ServerSideErrorResponse("Something went wrong while ending friendship, try again.", reqErr, "Error in friend req handler"), nil
		}

		return models.SuccessfulRequestResponse("Friendship ended", false), nil

	default:
		return models.ServerSideErrorResponse("Something went wrong while determining friendship action, try again.", fmt.Errorf("Status returned does not match any expected outcome. status returned: %v", status), "Status returned does not match any expected outcome"), nil
	}
}
