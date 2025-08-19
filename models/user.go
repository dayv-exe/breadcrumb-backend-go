package models

// standard user db model

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type User struct {
	Userid        string
	Nickname      string
	Name          string
	Bio           string
	DpUrl         string
	isSuspended   bool
	isDeactivated bool
	UserLogs      UserLogs
}

type CognitoInfo struct {
	Cognito map[string]string `json:"cognito"`
}

func NewUser(userid string, nickname string, name string, isSuspended bool) User {
	return User{
		Userid:        userid,
		Nickname:      nickname,
		Name:          name,
		Bio:           "",
		isSuspended:   isSuspended,
		isDeactivated: false,
		UserLogs:      NewUserLogs(),
	}
}

func (u User) DatabaseFormat() map[string]types.AttributeValue {
	logsJSON := u.UserLogs.DatabaseFormat()

	return map[string]types.AttributeValue{
		"pk":             &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", u.Userid)},
		"sk":             &types.AttributeValueMemberS{Value: "PROFILE"},
		"name":           &types.AttributeValueMemberS{Value: u.Name},
		"nickname":       &types.AttributeValueMemberS{Value: strings.ToLower(u.Nickname)},
		"bio":            &types.AttributeValueMemberS{Value: ""},
		"dpUrl":          &types.AttributeValueMemberS{Value: u.DpUrl},
		"is_suspended":   &types.AttributeValueMemberBOOL{Value: u.isSuspended},
		"is_deactivated": &types.AttributeValueMemberBOOL{Value: u.isDeactivated},
		"user_logs":      &types.AttributeValueMemberM{Value: logsJSON},
	}
}

func GetCognitoInfo(sub string, cognitoClient *cognitoidentityprovider.Client, ctx context.Context, userPoolId string) (*CognitoInfo, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(userPoolId),
		Username:   aws.String(sub),
	}

	user, err := getCognitoInfoQueryRunner(func() (*cognitoidentityprovider.AdminGetUserOutput, error) {
		return cognitoClient.AdminGetUser(ctx, input)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func getCognitoInfoQueryRunner(queryFn func() (*cognitoidentityprovider.AdminGetUserOutput, error)) (*CognitoInfo, error) {
	// THIS CUSTOM FUNCTION TO BE TESTED
	output, err := queryFn()

	if err != nil {
		return nil, err
	}

	userInfo := map[string]string{}

	for _, attr := range output.UserAttributes {
		userInfo[*attr.Name] = *attr.Value
	}

	return &CognitoInfo{
		Cognito: map[string]string{
			"email":     userInfo["email"],
			"birthdate": userInfo["birthdate"],
		},
	}, nil
}

func GetDbInfo(username string, tableName string, ctx context.Context, dbClient dynamodb.Client) (*User, error) {
	return nil, nil
}

func getDbInfoQueryRunner(queryFn func() (*dynamodb.QueryOutput, error)) (*User, error) {
	_, err := queryFn()

	if err != nil {
		return nil, err
	}

	user := User{}

	return &user, nil
}
