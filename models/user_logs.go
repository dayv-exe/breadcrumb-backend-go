package models

import (
	"breadcrumb-backend-go/utils"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserLogs struct {
	DateJoined               string
	BirthdateChangeCount     int
	LastNicknameChange       string
	LastEmailChange          string
	LastLogin                string
	forceChangeNickname      bool
	suspensionReason         string
	defaultProfilePicFgColor string
	defaultProfilePicBgColor string
}

func NewUserLogs() UserLogs {
	defaultColors := utils.GenerateRandomColorPair()
	return UserLogs{
		DateJoined:               utils.GetTimeNow(),
		BirthdateChangeCount:     0,
		LastLogin:                utils.GetTimeNow(),
		forceChangeNickname:      false,
		defaultProfilePicFgColor: defaultColors.Foreground,
		defaultProfilePicBgColor: defaultColors.Background,
	}
}

func (ul UserLogs) DatabaseFormat() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"date_joined":            &types.AttributeValueMemberS{Value: ul.DateJoined},
		"birthdate_change_count": &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", ul.BirthdateChangeCount)},
		"last_nickname_change":   &types.AttributeValueMemberS{Value: ul.LastNicknameChange},
		"last_email_change":      &types.AttributeValueMemberS{Value: ul.LastEmailChange},
		"last_login":             &types.AttributeValueMemberS{Value: ul.LastLogin},
		"force_change_nickname":  &types.AttributeValueMemberBOOL{Value: ul.forceChangeNickname},
		"default_pic_fg":         &types.AttributeValueMemberS{Value: ul.defaultProfilePicFgColor},
		"default_pic_bg":         &types.AttributeValueMemberS{Value: ul.defaultProfilePicBgColor},
		"suspension_reason":      &types.AttributeValueMemberS{Value: ul.suspensionReason},
	}
}
