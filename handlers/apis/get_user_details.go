package apis

import (
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"log"

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

func isAuthenticatedUser(req *events.APIGatewayProxyRequest, userId string) bool {
	sub := utils.GetAuthUserId(req)
	if sub == "" {
		return false
	}
	return sub == userId
}

func (deps *GetUserDetailsDependencies) HandleGetUserDetails(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// gets user details from db, returns only public information, if user requesting is not the user being requested.

	log.Println("Function begins")

	nickname, exists := req.PathParameters["nickname"]

	if !exists {
		log.Println("no nickname parsed")
		return models.InvalidRequestErrorResponse(""), nil
	}

	userDbHelper := models.UserDbHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       &ctx,
	}

	user, err := userDbHelper.FindByNickname(nickname)

	// error
	if err != nil {
		log.Println("ERROR OCCURRED: " + err.Error())
		return models.ServerSideErrorResponse("", err), nil
	}

	// no user found
	if user == nil {
		log.Println("no user found")
		return models.NotFoundResponse("User not found"), nil
	}

	if isAuthenticatedUser(req, user.Userid) {
		log.Println("checking if user is requesting their own profile")
		// if the logged in user is requesting their own information
		userCognitoHelper := models.UserCognitoHelper{
			UserPoolId:    deps.UserPoolId,
			CognitoClient: deps.CognitoClient,
			Ctx:           &ctx,
		}

		log.Println("Getting cognito info")
		userCognitoInfo, err := userCognitoHelper.GetCognitoInfo(user.Userid)

		if err != nil {
			log.Println("error")
			return models.ServerSideErrorResponse("Something went wrong.", err), nil
		}

		if userCognitoInfo == nil {
			log.Println("no cognito user found")
			return models.NotFoundResponse("User details not found."), nil
		}

		log.Println("user is requesting their own data")
		return models.SuccessfulGetRequestResponse(fullUserInfo{
			User:        *user,
			CognitoInfo: *userCognitoInfo,
		}), nil
	}

	log.Println("User is requesting other users data")
	// only return nickname, name, profile picture
	return models.SuccessfulGetRequestResponse(minUserInfo{
		Nickname: user.Nickname,
		Name:     user.Name,
		DpUrl:    user.DpUrl,
	}), nil
}
