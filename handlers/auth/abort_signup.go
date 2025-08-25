package auth

import (
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// deletes an unverified user from cognito when the user cancels verification process on the frontend

type AbortSignupDependencies struct {
	Client     *cognitoidentityprovider.Client
	UserPoolId string
}

func (deps *AbortSignupDependencies) AbortSignupHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// invalid username(sub, uuid)
	userId := req.PathParameters["id"]
	if userId == "" {
		return models.InvalidRequestErrorResponse(""), nil
	}

	cognitoHelper := helpers.UserCognitoHelper{
		UserPoolId:    deps.UserPoolId,
		CognitoClient: deps.Client,
		Ctx:           ctx,
	}

	err := cognitoHelper.DeleteFromCognito(userId, true)

	if err != nil {
		return models.ServerSideErrorResponse("An error occurred while trying to remove your account.", err, "error while trying to delete user from cognito"), nil
	}

	return models.SuccessfulRequestResponse("successfully cancelled signup.", false), nil
}
