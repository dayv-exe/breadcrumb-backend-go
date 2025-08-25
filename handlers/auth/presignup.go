package auth

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/utils"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PreSignupDependencies struct {
	DdbClient *dynamodb.Client
	TableName string
}

func (deps *PreSignupDependencies) PreSignupHandler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	// runs before user is added to cognito user pool

	nickname := event.Request.UserAttributes["nickname"]
	birthdate := event.Request.UserAttributes["birthdate"]
	fullname := event.Request.UserAttributes["name"]

	// nickname check
	if !utils.NicknameValid(nickname) {
		return event, fmt.Errorf("invalid nickname")
	}

	dbHelper := helpers.UserDynamoHelper{
		DbClient:  deps.DdbClient,
		TableName: deps.TableName,
		Ctx:       ctx,
	}

	// one final check if username is still free
	nicknameAvail, err := dbHelper.NicknameAvailable(nickname)

	if err != nil {
		return event, fmt.Errorf("error checking nickname availability %w", err)
	}

	if !nicknameAvail {
		return event, fmt.Errorf("nickname taken")
	}

	// birthdate check
	validBirthdate, err := utils.BirthdateIsValid(birthdate)

	if err != nil {
		return event, fmt.Errorf("Birthdate is in a wrong format, it should be DD/MM/YYYY! ERROR: %s.", err)
	}

	if !validBirthdate {
		return event, fmt.Errorf("Birthdate is invalid, users must be between %d and %d years old, expected format is dd/mm/yyyy", constants.MIN_AGE, constants.MAX_AGE)
	}

	// fullname check
	if !utils.NameIsValid(&fullname) {
		return event, fmt.Errorf("Fullname cannot be longer than %d characters", constants.MAX_FULLNAME_CHARS)
	}

	return event, nil
}
