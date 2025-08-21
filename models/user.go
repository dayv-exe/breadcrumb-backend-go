package models

// standard user db model

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type User struct {
	Userid        string   `dynamodbav:"pk" json:"userId"`
	DbDescription string   `dynamodbav:"sk"`
	Nickname      string   `dynamodbav:"nickname" json:"nickname"`
	Name          string   `dynamodbav:"name" json:"name"`
	Bio           string   `dynamodbav:"bio" json:"bio"`
	DpUrl         string   `dynamodbav:"dpUrl" json:"dpUrl"`
	IsSuspended   bool     `dynamodbav:"is_suspended" json:"isSuspended"`
	IsDeactivated bool     `dynamodbav:"is_deactivated" json:"isDeactivated"`
	UserLogs      UserLogs `dynamodbav:"user_logs" json:"userLogs"`
}

type CognitoInfo struct {
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
}

type UserDbHelper struct {
	DbClient  *dynamodb.Client
	TableName string
	Ctx       *context.Context
}

type UserCognitoHelper struct {
	UserPoolId    string
	CognitoClient *cognitoidentityprovider.Client
	Ctx           *context.Context
}

func NewUser(userid string, nickname string, name string, isSuspended bool) *User {
	return &User{
		Userid:        userid,
		Nickname:      nickname,
		Name:          name,
		Bio:           "",
		IsSuspended:   isSuspended,
		IsDeactivated: false,
		UserLogs:      *NewUserLogs(),
	}
}

func (deps UserDbHelper) AddNewUser(userid string, nickname string, name string, isSuspended bool) error {
	// create user
	user := NewUser(userid, nickname, name, false)

	nicknameItem := GetNicknameDbItem(user)

	if nicknameItem == nil {
		return fmt.Errorf("Unable to create nickname item.")
	}

	// add to db
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				// add user to db
				Put: &types.Put{
					TableName: aws.String(deps.TableName),
					Item:      user.DatabaseFormat(),
				},
			},

			{
				// adds nickname item to reserve name
				Put: &types.Put{
					TableName: aws.String(deps.TableName),
					Item:      nicknameItem,
					// if this fails most likely because nickname is already in use everything will roll back
					ConditionExpression: aws.String("attribute_not_exists(pk)"),
				},
			},
		},
	}

	_, err := deps.DbClient.TransactWriteItems(*deps.Ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func (u User) DatabaseFormat() map[string]types.AttributeValue {
	u.Userid = "USER#" + u.Userid
	u.DbDescription = "PROFILE"
	item, err := attributevalue.MarshalMap(u)

	if err != nil {
		return nil
	}
	return item
}

func (deps UserCognitoHelper) GetCognitoInfo(sub string) (*CognitoInfo, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(deps.UserPoolId),
		Username:   aws.String(sub),
	}

	user, err := getCognitoInfoQueryRunner(func() (*cognitoidentityprovider.AdminGetUserOutput, error) {
		return deps.CognitoClient.AdminGetUser(*deps.Ctx, input)
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
		Email:     userInfo["email"],
		Birthdate: userInfo["birthdate"],
	}, nil
}

func convertToUser(item map[string]types.AttributeValue) (*User, error) {
	var u User
	if err := attributevalue.UnmarshalMap(item, &u); err != nil {
		return nil, err
	}

	u.Userid = strings.TrimPrefix(u.Userid, "USER#")

	return &u, nil
}

func (deps UserDbHelper) FindByNickname(nickname string) (*User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(deps.TableName),
		IndexName:              aws.String("NicknameIndex"),
		KeyConditionExpression: aws.String("nickname = :nick"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nick": &types.AttributeValueMemberS{Value: nickname},
		},
		Limit: aws.Int32(1),
	}

	output, err := deps.DbClient.Query(*deps.Ctx, input)

	if err != nil {
		return nil, err
	}

	if output.Count < 1 {
		return nil, nil
	}

	return convertToUser(output.Items[0])
}
