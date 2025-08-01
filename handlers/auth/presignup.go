package auth

import (
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

	// nickname check
	if !utils.IsNicknameValid(nickname) {
		return event, fmt.Errorf("invalid nickname")
	}

	nicknameAvail, err := utils.IsNicknameAvailableInDynamodb(nickname, deps.TableName, deps.DdbClient, ctx)

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
		return event, fmt.Errorf("Birthdate is invalid, users must be between 13 and 85 years old, expected format is dd/mm/yyyy")
	}

	return event, nil
}
