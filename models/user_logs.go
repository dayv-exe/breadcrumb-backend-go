package models

import (
	"breadcrumb-backend-go/utils"
	"encoding/json"
	"log"

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

func (ul *UserLogs) ToMap() map[string]interface{} {
	data, err := json.Marshal(ul)
	if err != nil {
		log.Fatalf("Error marshaling UserLogs: %v", err)
		return nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalf("Error unmarshaling UserLogs to map: %v", err)
		return nil
	}

	return result
}

func ToStruct(m map[string]interface{}) (*UserLogs, error) {
	// Convert map to JSON
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// Convert JSON to struct
	var ul UserLogs
	if err := json.Unmarshal(data, &ul); err != nil {
		return nil, err
	}

	return &ul, nil
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
