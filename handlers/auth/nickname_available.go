package auth

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type NicknameDependencies struct {
	DdbClient *dynamodb.Client
	TableName string
}

func (deps *NicknameDependencies) HandleNicknameAvailable(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// returns true if a nickname has not been taken by a user in dynamodb

	// if nickname is invalid return false immediately
	nickname := strings.ToLower(req.PathParameters["nickname"])
	if nickname == "" || !utils.IsNicknameValid(nickname) {
		return models.SuccessfulRequestResponse(fmt.Sprintf("%v", false), false), nil
	}

	isAvailable, dbErr := utils.IsNicknameAvailableInDynamodb(nickname, deps.TableName, deps.DdbClient, ctx)

	if dbErr != nil {
		return models.ServerSideErrorResponse("An error has occurred, try again.", dbErr), nil
	}

	return models.SuccessfulRequestResponse(fmt.Sprintf("%v", isAvailable), false), nil
}
