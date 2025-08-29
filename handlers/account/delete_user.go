package account

import (
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DeleteUserDependencies struct {
	DbClient  *dynamodb.Client
	TableName string

	CognitoClient *cognitoidentityprovider.Client
	UserPoolId    string
}

func (deps *DeleteUserDependencies) HandleDeleteUser(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId := utils.GetAuthUserId(req)

	// get user details from db
	dbHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}
	user, uErr := dbHelper.FindById(userId)

	if uErr != nil {
		return models.ServerSideErrorResponse("An error occurred while trying to delete your account, try again", uErr, "Error from find by id"), nil
	}

	if user == nil {
		return models.NotFoundResponse(""), nil
	}

	// delete user from dynamodb
	delErr := dbHelper.DeleteFromDynamo(userId, user.Nickname)

	if delErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while trying to delete your account, try again", delErr, "error from delete from dynamo db"), nil
	}

	// delete user from cognito
	cognitoHelper := helpers.UserCognitoHelper{
		CognitoClient: deps.CognitoClient,
		UserPoolId:    deps.UserPoolId,
		Ctx:           ctx,
	}

	cogErr := cognitoHelper.DeleteFromCognito(userId, true)

	if cogErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while trying to delete your account, try again.", cogErr, "error from delete from cognito"), nil
	}

	return models.SuccessfulRequestResponse("", false), nil
}
