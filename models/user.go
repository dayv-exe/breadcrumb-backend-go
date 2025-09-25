package models

// standard user db model

import (
	"breadcrumb-backend-go/utils"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Userid                  string `dynamodbav:"pk" json:"userId"`
	DbDescription           string `dynamodbav:"sk" json:"type"`
	Nickname                string `dynamodbav:"nickname" json:"nickname"`
	Name                    string `dynamodbav:"name" json:"name"`
	Bio                     string `dynamodbav:"bio" json:"bio"`
	DpUrl                   string `dynamodbav:"dpUrl" json:"dpUrl"`
	IsSuspended             bool   `dynamodbav:"is_suspended" json:"isSuspended"`
	IsDeactivated           bool   `dynamodbav:"is_deactivated" json:"isDeactivated"`
	DateJoined              string `dynamodbav:"date_joined" json:"dateJoined"`
	BirthdateChangeCount    int    `dynamodbav:"birthdate_change_count" json:"birthdateChangeCount"`
	LastNicknameChange      string `dynamodbav:"last_nickname_change" json:"lastNicknameChange"`
	LastEmailChange         string `dynamodbav:"last_email_change" json:"lastEmailChange"`
	LastLogin               string `dynamodbav:"last_login" json:"lastLogin"`
	LastNameChange          string `dynamodb:"last_name_change" json:"lastNameChange"`
	ForceChangeNickname     bool   `dynamodbav:"force_change_nickname" json:"forceChangeNickname"`
	SuspensionReason        string `dynamodbav:"suspension_reason" json:"suspensionReason"`
	DefaultProfilePicColors string `dynamodbav:"default_pic_colors" json:"defaultPicColors"`
}

type PrimaryUserInfo struct {
	Userid                  string `json:"userId"`
	Nickname                string `json:"nickname"`
	Name                    string `json:"name"`
	DpUrl                   string `json:"dpUrl"`
	Bio                     string `json:"bio"`
	DefaultProfilePicColors string `json:"defaultPicColors"`
	IsSuspended             bool   `json:"isSuspended"`
	IsDeactivated           bool   `json:"isDeactivated"`
	FriendshipStatus        string `json:"friends"`
}

const (
	UserPkPrefix = "USER#"
	UserSkPrefix = "PROFILE"
)

func NewPrimaryUserInfo(u *User, friendshipStatus string) *PrimaryUserInfo {
	return &PrimaryUserInfo{
		Userid:                  u.Userid,
		Nickname:                u.Nickname,
		Name:                    u.Name,
		DpUrl:                   u.DpUrl,
		Bio:                     u.Bio,
		DefaultProfilePicColors: u.DefaultProfilePicColors,
		IsSuspended:             u.IsSuspended,
		IsDeactivated:           u.IsDeactivated,
		FriendshipStatus:        friendshipStatus,
	}
}

func NewUser(userid string, nickname string, name string, isSuspended bool) *User {
	defaultColors := utils.GenerateRandomColorPair()

	return &User{
		Userid:                  userid,
		Nickname:                strings.ToLower(nickname),
		Name:                    name,
		Bio:                     "",
		IsSuspended:             isSuspended,
		IsDeactivated:           false,
		DateJoined:              utils.GetTimeNow(),
		BirthdateChangeCount:    0,
		LastNicknameChange:      "",
		LastEmailChange:         "",
		LastNameChange:          "",
		LastLogin:               utils.GetTimeNow(),
		ForceChangeNickname:     false,
		SuspensionReason:        "",
		DefaultProfilePicColors: defaultColors.Foreground + utils.AddPrefix("", defaultColors.Background),
		DbDescription:           UserSkPrefix,
	}
}

func UserKey(userid string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: utils.AddPrefix(UserPkPrefix, userid)},
		"sk": &types.AttributeValueMemberS{Value: UserSkPrefix},
	}
}

func (u *User) DatabaseFormat() *map[string]types.AttributeValue {
	u.Userid = utils.AddPrefix(UserPkPrefix, u.Userid)

	item, err := attributevalue.MarshalMap(u)

	if err != nil {
		return nil
	}
	return &item
}

func ConvertToUser(item map[string]types.AttributeValue) (*User, error) {
	var u User
	if err := attributevalue.UnmarshalMap(item, &u); err != nil {
		return nil, err
	}

	u.Userid = strings.TrimPrefix(u.Userid, UserPkPrefix)

	return &u, nil
}
