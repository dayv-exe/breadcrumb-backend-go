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

	// if nickname is invalid return false immediately
	nickname := strings.ToLower(req.PathParameters["nickname"])
	if nickname == "" || !utils.IsNicknameValid(nickname) {
		return models.SuccessfulRequestResponse(fmt.Sprintf("%v", false), false), nil
	}

	takenInDb, dbErr := utils.IsNicknameTakenInDynamodb(nickname, deps.TableName, deps.DdbClient, ctx)

	if dbErr != nil {
		return models.ServerSideErrorResponse("An error has occurred, try again.", dbErr.Error()), nil
	}

	available := !takenInDb

	return models.SuccessfulRequestResponse(fmt.Sprintf("%v", available), false), nil
}
