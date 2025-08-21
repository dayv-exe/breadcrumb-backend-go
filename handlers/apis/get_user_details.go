package apis

import (
	"breadcrumb-backend-go/models"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// if a user request their details, this function will return all their info from cognito and dynamodb
// if a user request details of another user, it will only return profile picture url, nickname, username, and maybe mutual friends

type UserDetailsDependencies struct {
	// cognito stuff
	CognitoClient *cognitoidentityprovider.Client
	UserPoolId    string

	// dynamodb stuff
	DdbClient *dynamodb.Client
	TableName string
}

type minUserInfo struct {
	// the least amount of information on any user
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	DpUrl    string `json:"dpUrl"`
}

type fullUserInfo struct {
	models.User
	models.CognitoInfo
}

func getAuthUserId(req *events.APIGatewayProxyRequest) string {
	claims, ok := req.RequestContext.Authorizer["claims"].(map[string]interface{})
	if !ok {
		// missing claims
		return ""
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		// missing sub
		return ""
	}

	return sub
}

func isAuthenticatedUser(req *events.APIGatewayProxyRequest, userId string) bool {
	return getAuthUserId(req) == userId
}

func (deps *UserDetailsDependencies) HandleOtherUserDetails(ctx *context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// gets user details from db, returns only public information, if user requesting is not the user being requested.

	nickname, exists := req.PathParameters["nickname"]

	if !exists {
		return models.InvalidRequestErrorResponse(""), nil
	}

	user, err := models.FindByNickname(nickname, deps.TableName, ctx, deps.DdbClient)

	if err != nil {
		return models.ServerSideErrorResponse("", err), nil
	}

	if user == nil {
		return models.NotFoundResponse("User not found"), nil
	}

	if isAuthenticatedUser(req, user.Userid) {
		// if the logged in user is requesting their own information
		userCognitoInfo, err := models.GetCognitoInfo(user.Userid, deps.CognitoClient, ctx, deps.UserPoolId)

		if err != nil {
			return models.ServerSideErrorResponse("Something went wrong.", err), nil
		}

		if userCognitoInfo == nil {
			return models.NotFoundResponse("User details not found."), nil
		}

		return models.SuccessfulGetRequestResponse(fullUserInfo{
			User:        *user,
			CognitoInfo: *userCognitoInfo,
		}), nil
	}

	// only return nickname, name, profile picture
	return models.SuccessfulGetRequestResponse(minUserInfo{
		Nickname: user.Nickname,
		Name:     user.Name,
		DpUrl:    user.DpUrl,
	}), nil
}
