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
	DdbClient       *dynamodb.Client
	TableName       string
	SearchTableName string
	CognitoClient   *cognitoidentityprovider.Client
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
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	// the add users names to search table
	searchHelper := helpers.SearchDynamoHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.SearchTableName,
		Ctx:       ctx,
	}

	// create new user
	newUser := models.NewUser(userID, nickName, name, false)

	err := dbHelper.AddUser(newUser)                     // adds new user to users table
	indexErr := searchHelper.AddUserSearchIndex(newUser) // adds users names to search table

	if err != nil || indexErr != nil {
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

		if err == nil && indexErr != nil {
			// if we are unable to add the users names(nickname and full name) to the search table, we will delete the user from dynamo since they must have already been added to users table at this point
			dynamodbErr := dbHelper.DeleteFromDynamo(userID, nickName)

			if dynamodbErr != nil {
				log.Println("Error occurred while trying to remove user dynamodb details: " + dynamodbErr.Error())
				return nil, fmt.Errorf("Something went wrong while setting up account, try again.")
			}
		} else if err != nil && indexErr == nil {
			// if we are unable to add the user to dynamodb, then delete the users names index from search table
			delIndexErr := searchHelper.DeleteUserIndexes(newUser)
			if delIndexErr != nil {
				log.Println("Error occurred while trying to remove user search indexes: " + delIndexErr.Error())
				return nil, fmt.Errorf("Something went wrong while setting up account, try again.")
			}
		}

		return nil, fmt.Errorf("Something went wrong while creating new account, try again.")
	}

	return event, nil
}
