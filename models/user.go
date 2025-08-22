package models

// standard user db model

import (
	"breadcrumb-backend-go/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type User struct {
	Userid                   string `dynamodbav:"pk" json:"userId"`
	DbDescription            string `dynamodbav:"sk" json:"type"`
	Nickname                 string `dynamodbav:"nickname" json:"nickname"`
	Name                     string `dynamodbav:"name" json:"name"`
	Bio                      string `dynamodbav:"bio" json:"bio"`
	DpUrl                    string `dynamodbav:"dpUrl" json:"dpUrl"`
	IsSuspended              bool   `dynamodbav:"is_suspended" json:"isSuspended"`
	IsDeactivated            bool   `dynamodbav:"is_deactivated" json:"isDeactivated"`
	DateJoined               string `dynamodbav:"date_joined" json:"dateJoined"`
	BirthdateChangeCount     int    `dynamodbav:"birthdate_change_count" json:"birthdateChangeCount"`
	LastNicknameChange       string `dynamodbav:"last_nickname_change" json:"lastNicknameChange"`
	LastEmailChange          string `dynamodbav:"last_email_change" json:"lastEmailChange"`
	LastLogin                string `dynamodbav:"last_login" json:"lastLogin"`
	ForceChangeNickname      bool   `dynamodbav:"force_change_nickname" json:"forceChangeNickname"`
	SuspensionReason         string `dynamodbav:"suspension_reason" json:"suspensionReason"`
	DefaultProfilePicFgColor string `dynamodbav:"default_pic_fg" json:"defaultPicFg"`
	DefaultProfilePicBgColor string `dynamodbav:"default_pic_bg" json:"defaultPicBg"`
}

type MinUserInfo struct {
	// the least amount of information on any user
	Nickname                 string `json:"nickname"`
	Name                     string `json:"name"`
	DpUrl                    string `json:"dpUrl"`
	DefaultProfilePicFgColor string `json:"defaultPicFg"`
	DefaultProfilePicBgColor string `json:"defaultPicBg"`
	IsSuspended              bool   `json:"isSuspended"`
	IsDeactivated            bool   `json:"isDeactivated"`
}

type FullUserInfo struct {
	// most amount of information to be returned on a user
	User
	CognitoInfo
}

type CognitoInfo struct {
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
}

type UserDbHelper struct {
	DbClient  *dynamodb.Client
	TableName string
	Ctx       context.Context
}

type UserCognitoHelper struct {
	UserPoolId    string
	CognitoClient *cognitoidentityprovider.Client
	Ctx           context.Context
}

func NewUser(userid string, nickname string, name string, isSuspended bool) *User {
	defaultColors := utils.GenerateRandomColorPair()

	return &User{
		Userid:                   userid,
		Nickname:                 strings.ToLower(nickname),
		Name:                     name,
		Bio:                      "",
		IsSuspended:              isSuspended,
		IsDeactivated:            false,
		DateJoined:               utils.GetTimeNow(),
		BirthdateChangeCount:     0,
		LastNicknameChange:       "",
		LastEmailChange:          "",
		LastLogin:                utils.GetTimeNow(),
		ForceChangeNickname:      false,
		SuspensionReason:         "",
		DefaultProfilePicFgColor: defaultColors.Foreground,
		DefaultProfilePicBgColor: defaultColors.Background,
	}
}

func (deps *UserDbHelper) AddNewUser(userid string, nickname string, name string, isSuspended bool) error {
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
				// adds nickname item to reserve name
				Put: &types.Put{
					TableName: aws.String(deps.TableName),
					Item:      nicknameItem,
					// if this fails most likely because nickname is already in use everything will roll back
					ConditionExpression: aws.String("attribute_not_exists(pk)"),
				},
			},

			{
				// add user to db
				Put: &types.Put{
					TableName: aws.String(deps.TableName),
					Item:      user.databaseFormat(),
				},
			},
		},
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) databaseFormat() map[string]types.AttributeValue {
	u.Userid = "USER#" + u.Userid
	u.DbDescription = "PROFILE"
	item, err := attributevalue.MarshalMap(u)

	if err != nil {
		return nil
	}
	return item
}

func (deps *UserCognitoHelper) GetCognitoInfo(sub string) (*CognitoInfo, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(deps.UserPoolId),
		Username:   aws.String(sub),
	}

	user, err := getCognitoInfoQueryRunner(func() (*cognitoidentityprovider.AdminGetUserOutput, error) {
		return deps.CognitoClient.AdminGetUser(deps.Ctx, input)
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
	log.Println("Started converting user")
	if err := attributevalue.UnmarshalMap(item, &u); err != nil {
		return nil, err
	}

	u.Userid = strings.TrimPrefix(u.Userid, "USER#")

	log.Println("done converting user")
	return &u, nil
}

func (deps *UserDbHelper) FindByNickname(nickname string) (*User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(deps.TableName),
		IndexName:              aws.String("NicknameIndex"),
		KeyConditionExpression: aws.String("nickname = :nick"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nick": &types.AttributeValueMemberS{Value: strings.ToLower(nickname)},
		},
		Limit: aws.Int32(1),
	}

	output, err := deps.DbClient.Query(deps.Ctx, input)

	if err != nil {
		return nil, err
	}

	if output.Count < 1 {
		return nil, nil
	}

	log.Println("Done fetching user, got %w", output.Items[0])

	return convertToUser(output.Items[0])
}

func (deps *UserDbHelper) FindById(id string) (*User, error) {
	input := dynamodb.GetItemInput{
		Key:       profileKey(id),
		TableName: &deps.TableName,
	}

	output, err := deps.DbClient.GetItem(deps.Ctx, &input)

	if err != nil {
		return nil, err
	}

	return convertToUser(output.Item)
}

func nicknameKey(nickname string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "NICKNAME#" + nickname},
		"sk": &types.AttributeValueMemberS{Value: "NICKNAME"},
	}
}

func profileKey(userid string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "USER#" + userid},
		"sk": &types.AttributeValueMemberS{Value: "PROFILE"},
	}
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
			var notFoundErr *cognitoTypes.UserNotFoundException
			if !errors.As(getUserErr, &notFoundErr) {
				return getUserErr
			}
		}

		// end delete process without deleting user if the user is anything but unconfirmed
		// only deletes unconfirmed user if the ignore confirmation flag is set to true
		if getUserOutput.UserStatus != cognitoTypes.UserStatusTypeUnconfirmed {
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
		var notFoundErr *cognitoTypes.UserNotFoundException
		if errors.As(err, &notFoundErr) {
			log.Println("Could not find user account to delete: " + id)
			return nil // simulate successful deletion
		}
		return fmt.Errorf("error occurred while trying to delete a user %w", err)
	}

	return nil
}

func (deps *UserDbHelper) DeleteFromDynamo(userId string, nickname string) error {
	// delete user profile, nickname, friends, post and allat
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				// delete account metadata
				Delete: &types.Delete{
					Key:       profileKey(userId),
					TableName: aws.String(deps.TableName),
				},
			},
			{
				// delete nickname reservation
				Delete: &types.Delete{
					Key:       nicknameKey(nickname),
					TableName: aws.String(deps.TableName),
				},
			},
		},
	}

	_, err := deps.DbClient.TransactWriteItems(deps.Ctx, input)

	if err != nil {
		var notFoundErr *types.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			log.Println("Could not find user resource to delete: %w", nickname)
			return nil
		}
		return err
	}

	return nil
}
