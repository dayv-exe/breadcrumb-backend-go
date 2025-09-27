package auth

import (
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PostConfirmationDependencies struct {
	DdbClient     *dynamodb.Client
	CognitoClient *cognitoidentityprovider.Client
	TableNames    *utils.TableNames
}

func (deps PostConfirmationDependencies) HandlePostConfirmation(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (interface{}, error) {
	// runs after user has validated their email

	// this function should only run if it is trigger by signup confirm
	if event.TriggerSource != "PostConfirmation_ConfirmSignUp" {
		return event, nil
	}

	userID := event.Request.UserAttributes["sub"]
	nickName := event.Request.UserAttributes["nickname"]
	name := event.Request.UserAttributes["name"]

	// to add user to users table
	dbHelper := helpers.UserDynamoHelper{
		DbClient:   deps.DdbClient,
		TableNames: deps.TableNames,
		Ctx:        ctx,
	}

	// create new user
	newUser := models.NewUser(userID, nickName, name, false)

	err := dbHelper.AddUser(newUser) // adds new user to users table

	if err != nil {
		// if something goes wrong during the signup process deelete user cognito info
		log.Printf("ERROR IN SIGNUP CONFIRM GO FUNC: %v", err)

		// remove the users info from cognito
		cognitoHelper := helpers.UserCognitoHelper{
			UserPoolId:    event.UserPoolID,
			CognitoClient: deps.CognitoClient,
			Ctx:           ctx,
		}
		cognitoErr := cognitoHelper.DeleteFromCognito(userID, true)
		if cognitoErr != nil {
			log.Println("Error occurred while trying to remove user cognito account: " + cognitoErr.Error())
			return nil, fmt.Errorf("Something went wrong while setting up account, try again.")
		}
		log.Print("DELETED FROM COGNITO")

		return nil, fmt.Errorf("Something went wrong while creating new account, try again.")
	}

	return event, nil
}
