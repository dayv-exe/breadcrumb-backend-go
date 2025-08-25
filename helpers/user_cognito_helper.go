package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type UserCognitoHelper struct {
	CognitoClient *cognitoidentityprovider.Client
	UserPoolId    string
	Ctx           context.Context
}

type CognitoManagedInfo struct {
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
}

func (deps *UserCognitoHelper) GetManagedInfo(sub string) (*CognitoManagedInfo, error) {
	// returns user details managed by cognito like email and birthdate

	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(deps.UserPoolId),
		Username:   aws.String(sub),
	}

	output, err := getManagedInfoCustomFilter(func() (*cognitoidentityprovider.AdminGetUserOutput, error) {
		return deps.CognitoClient.AdminGetUser(deps.Ctx, input)
	})

	if err != nil {
		return nil, err
	}
	return output, nil
}

func getManagedInfoCustomFilter(queryFn func() (*cognitoidentityprovider.AdminGetUserOutput, error)) (*CognitoManagedInfo, error) {
	output, err := queryFn()

	if err != nil {
		return nil, err
	}

	userAttr := map[string]string{}

	for _, attr := range output.UserAttributes {
		userAttr[*attr.Name] = *attr.Value
	}

	return &CognitoManagedInfo{
		Email:     userAttr["email"],
		Birthdate: userAttr["birthdate"],
	}, nil
}

func (deps *UserCognitoHelper) DeleteFromCognito(id string, ignoreConfirmationStatus bool) error {

	// deletes unconfirmed users, except ignore confirmation status then it deletes any user

	if !ignoreConfirmationStatus {
		getUserInput := &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(deps.UserPoolId),
			Username:   aws.String(id),
		}

		getUserOutput, getUserErr := deps.CognitoClient.AdminGetUser(deps.Ctx, getUserInput)

		// return error only if the error is not a user not found exception
		if getUserErr != nil {
			var notFoundErr *types.UserNotFoundException
			if errors.As(getUserErr, &notFoundErr) {
				return nil
			}

			return getUserErr
		}

		// end delete process without deleting user if the user is anything but unconfirmed
		// only deletes unconfirmed user if the ignore confirmation flag is set to true
		if getUserOutput.UserStatus != types.UserStatusTypeUnconfirmed {
			log.Println("Failed to delete user because they have been confirmed and function is not ignoring confirmation status")
			return nil
		}
	}

	input := &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(deps.UserPoolId),
		Username:   aws.String(id),
	}

	_, err := deps.CognitoClient.AdminDeleteUser(deps.Ctx, input)

	if err != nil {
		var notFoundErr *types.UserNotFoundException
		if errors.As(err, &notFoundErr) {
			log.Println("Could not find user account to delete: " + id)
			return nil // simulate successful deletion
		}
		return fmt.Errorf("error occurred while trying to delete a user %w", err)
	}

	return nil
}
