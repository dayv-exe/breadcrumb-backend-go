package auth

import (
	"breadcrumb-backend-go/models"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PostConfirmationDependencies struct {
	DdbClient *dynamodb.Client
	TableName string
}

func (deps PostConfirmationDependencies) HandlePostConfirmation(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (interface{}, error) {
	// runs after user has validated their email

	userID := event.Request.UserAttributes["sub"]
	nickName := event.Request.UserAttributes["nickname"]
	name := event.Request.UserAttributes["name"]

	// create new user
	user := models.UserDbHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       &ctx,
	}

	err := user.AddNewUser(userID, nickName, name, false)

	if err != nil {
		log.Println("ERROR IN SIGNUP CONFIRM GO FUNC: " + err.Error())
		return nil, fmt.Errorf("Something went wrong while creating new account, try again.")
	}

	return event, nil
}
