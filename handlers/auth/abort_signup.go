package auth

import (
	"breadcrumb-backend-go/models"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// deletes an unverified user from cognito when the user cancels verification process on the frontend

type AbortSignupDependencies struct {
	Client     *cognitoidentityprovider.Client
	UserPoolId string
}

func (asd *AbortSignupDependencies) AbortSignupHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// invalid username(sub, uuid)
	userId := req.PathParameters["id"]
	if userId == "" {
		return models.InvalidRequestErrorResponse(""), nil
	}

	// if user is verified then sign up has been completed and they thus cannot be deleted
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(asd.UserPoolId),
		Username:   aws.String(userId),
	}
	user, uErr := asd.Client.AdminGetUser(ctx, getUserInput)
	if uErr != nil {
		// if we cannot get user, most likely user id is invalid
		return models.InvalidRequestErrorResponse("invalid request"), nil
	}

	if user.UserStatus == types.UserStatusTypeConfirmed {
		// if user has been confirmed their signup cannot be aborted
		return models.InvalidRequestErrorResponse("Invalid request"), nil
	}

	input := &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(asd.UserPoolId),
		Username:   aws.String(userId),
	}

	_, err := asd.Client.AdminDeleteUser(ctx, input)

	if err != nil {
		return models.ServerSideErrorResponse("An error occurred while trying to delete the user", fmt.Errorf("error occurred while trying to delete a user %w", err)), err
	}

	return models.SuccessfulRequestResponse("", false), nil
}
