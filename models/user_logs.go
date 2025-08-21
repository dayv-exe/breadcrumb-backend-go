package models

import (
	"breadcrumb-backend-go/utils"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserLogs struct {
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

func NewUserLogs() UserLogs {
	defaultColors := utils.GenerateRandomColorPair()
	return UserLogs{
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

func (ul UserLogs) DatabaseFormat() map[string]types.AttributeValue {
	logs, err := attributevalue.MarshalMap(ul)

	if err != nil {
		return nil
	}

	return logs
}
