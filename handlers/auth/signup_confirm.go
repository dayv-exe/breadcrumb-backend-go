package auth

import (
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PostConfirmationDependencies struct {
	DdbClient     *dynamodb.Client
	TableName     string
	CognitoClient *cognitoidentityprovider.Client
}

func (deps PostConfirmationDependencies) HandlePostConfirmation(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (interface{}, error) {
	// runs after user has validated their email

	userID := event.Request.UserAttributes["sub"]
	nickName := event.Request.UserAttributes["nickname"]
	name := event.Request.UserAttributes["name"]

	dbHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	// create new user
	newUser := models.NewUser(userID, nickName, name, false)

	err := dbHelper.AddUser(newUser)

	if err != nil {
		// if somthing goes wrong during the signup process deelete user cognito info
		log.Println("ERROR IN SIGNUP CONFIRM GO FUNC: " + err.Error())

		cognitoHelper := helpers.UserCognitoHelper{
			UserPoolId:    event.UserPoolID,
			CognitoClient: deps.CognitoClient,
			Ctx:           ctx,
		}
		cognitoErr := cognitoHelper.DeleteFromCognito(userID, true)
		if cognitoErr != nil {
			log.Println("Error occurred while trying to remove user cognito account: " + cognitoErr.Error())
			return nil, fmt.Errorf("Something went wrong while creating new account, try again.")
		}

		return nil, fmt.Errorf("Something went wrong while creating new account, try again.")
	}

	return event, nil
}
