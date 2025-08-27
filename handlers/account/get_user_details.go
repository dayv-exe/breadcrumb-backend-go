package account

import (
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// if a user request their details, this function will return all their info from cognito and dynamodb
// if a user request details of another user, it will only return profile picture url, nickname, username, and maybe mutual friends

type GetUserDetailsDependencies struct {
	// cognito stuff
	CognitoClient *cognitoidentityprovider.Client
	UserPoolId    string

	// dynamodb stuff
	DdbClient *dynamodb.Client
	TableName string
}

type completeUserDetails struct {
	// all the information on a user
	models.User
	helpers.CognitoManagedInfo
}

func (deps *GetUserDetailsDependencies) HandleGetUserDetails(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// gets user details from db, returns only public information, if user requesting is not the user being requested.

	nickname, exists := req.PathParameters["nickname"]

	if !exists {
		return models.InvalidRequestErrorResponse(""), nil
	}

	// get all info on a user from dynamodb
	dbHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	user, dbErr := dbHelper.FindByNickname(nickname)

	// error
	if dbErr != nil {
		return models.ServerSideErrorResponse("", fmt.Errorf("Find by nickname error: %w", dbErr), "error while trying to find user by nickname"), nil
	}

	// no user found
	if user == nil {
		return models.NotFoundResponse("User not found"), nil
	}

	if utils.IsAuthenticatedUser(req, user.Userid) {
		// if the logged in user is requesting their own information
		cognitoHelper := helpers.UserCognitoHelper{
			CognitoClient: deps.CognitoClient,
			UserPoolId:    deps.UserPoolId,
			Ctx:           ctx,
		}

		userCognitoInfo, cogErr := cognitoHelper.GetManagedInfo(user.Userid)

		if cogErr != nil {
			return models.ServerSideErrorResponse("", fmt.Errorf("Get cognito info error: %w", cogErr), "while trying to get users cognito info"), nil
		}

		if userCognitoInfo == nil {
			return models.NotFoundResponse("User details not found."), nil
		}

		return models.SuccessfulGetRequestResponse(completeUserDetails{
			// return all the users info which is everything in dynamo and somethings in cognito
			*user,
			*userCognitoInfo,
		}), nil
	}

	// only return nickname, name, profile picture if one user requests another users information
	return models.SuccessfulGetRequestResponse(models.User{
		Nickname:                 user.Nickname,
		Name:                     user.Name,
		DpUrl:                    user.DpUrl,
		DefaultProfilePicFgColor: user.DefaultProfilePicFgColor,
		DefaultProfilePicBgColor: user.DefaultProfilePicBgColor,
		IsSuspended:              user.IsSuspended,
		IsDeactivated:            user.IsDeactivated,
	}), nil
}
