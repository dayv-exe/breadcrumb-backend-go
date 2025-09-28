package auth

import (
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type HandleNicknameAvailableDependencies struct {
	DdbClient *dynamodb.Client
	TableName string
}

func (deps *HandleNicknameAvailableDependencies) HandleNicknameAvailable(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// returns true if a nickname has not been taken by a user in dynamodb

	// if nickname is invalid return false immediately
	nickname := strings.ToLower(req.PathParameters["nickname"])
	if nickname == "" || !utils.NicknameValid(nickname) {
		return models.SuccessfulRequestResponse(fmt.Sprintf("%v", false), false), nil
	}

	dbHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DdbClient,
		Ctx:       ctx,
		TableName: deps.TableName,
	}

	isAvailable, dbErr := dbHelper.NicknameAvailable(nickname)

	if dbErr != nil {
		return models.ServerSideErrorResponse("An error has occurred, try again.", dbErr, "error while trying to check if nickname is available"), nil
	}

	return models.SuccessfulRequestResponse(fmt.Sprintf("%v", isAvailable), false), nil
}
