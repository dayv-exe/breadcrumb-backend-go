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

func getUser(useId bool, identifier string, dbHelper helpers.UserDynamoHelper) (*models.User, error) {
	if useId {
		return dbHelper.FindById(identifier)
	} else {
		return dbHelper.FindByNickname(identifier)
	}
}

func (deps *GetUserDetailsDependencies) HandleGetUserDetails(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// gets user details from db, returns only public information, if user requesting is not the user being requested.

	identifierName, idNameExists := req.PathParameters["identifier_name"]
	identifier, idExists := req.PathParameters["identifier"]

	if !idNameExists || !idExists {
		return models.InvalidRequestErrorResponse(""), nil
	}

	if identifierName != "nickname" && identifierName != "id" {
		return models.InvalidRequestErrorResponse("Invalid identifier name!"), nil
	}

	usingId := identifierName == "id"

	// get all info on a user from dynamodb
	dbHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	user, dbErr := getUser(usingId, identifier, dbHelper)

	// error
	if dbErr != nil {
		if usingId {
			return models.ServerSideErrorResponse("", fmt.Errorf("Find by id error: %w", dbErr), "error while trying to find user by id"), nil
		}
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
	friendshipHelper := helpers.FriendshipDynamoHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	friendshipStatus, fsErr := friendshipHelper.GetFriendshipStatus(utils.GetAuthUserId(req), user.Userid)

	if fsErr != nil {
		return models.ServerSideErrorResponse("Something went wrong while trying to get friendship status", fsErr, "error while trying to get friendship status"), nil
	}

	type resJson struct {
		models.User
		friends string `json:"friends"`
	}

	return models.SuccessfulGetRequestResponse(resJson{
		models.User{
			Userid:                   user.Userid,
			Nickname:                 user.Nickname,
			Name:                     user.Name,
			DpUrl:                    user.DpUrl,
			Bio:                      user.Bio,
			DefaultProfilePicFgColor: user.DefaultProfilePicFgColor,
			DefaultProfilePicBgColor: user.DefaultProfilePicBgColor,
			IsSuspended:              user.IsSuspended,
			IsDeactivated:            user.IsDeactivated,
		},
		friendshipStatus,
	}), nil
}
