package models

// standard user db model

import (
	"breadcrumb-backend-go/utils"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func (u *User) DatabaseFormat() map[string]types.AttributeValue {
	u.Userid = "USER#" + u.Userid
	u.DbDescription = "PROFILE"
	item, err := attributevalue.MarshalMap(u)

	if err != nil {
		return nil
	}
	return item
}

func ConvertToUser(item map[string]types.AttributeValue) (*User, error) {
	var u User
	log.Println("Started converting user")
	if err := attributevalue.UnmarshalMap(item, &u); err != nil {
		return nil, err
	}

	u.Userid = strings.TrimPrefix(u.Userid, "USER#")

	log.Println("done converting user")
	return &u, nil
}
