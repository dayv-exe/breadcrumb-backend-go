package models

// standard user db model

import (
	"breadcrumb-backend-go/utils"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	UserDisplayInfo
	UserAccountInfo
	UserPersonalInfo
	DbDescription string `dynamodbav:"sk" json:"type"`
}

type UserDisplayInfo struct {
	Userid                  string `dynamodbav:"pk" json:"userId"`
	Nickname                string `dynamodbav:"nickname" json:"nickname"`
	Name                    string `dynamodbav:"name" json:"name"`
	DpUrl                   string `dynamodbav:"dpUrl" json:"dpUrl"`
	DefaultProfilePicColors string `dynamodbav:"default_pic_colors" json:"defaultPicColors"`
}

type UserDisplayInfoNoId struct {
	Nickname                string `dynamodbav:"nickname" json:"nickname"`
	Name                    string `dynamodbav:"name" json:"name"`
	DpUrl                   string `dynamodbav:"dpUrl" json:"dpUrl"`
	DefaultProfilePicColors string `dynamodbav:"default_pic_colors" json:"defaultPicColors"`
}

type UserAccountInfo struct {
	Bio              string `dynamodbav:"bio" json:"bio"`
	FriendshipStatus string `dynamodbav:"friends" json:"friends"`
	IsSuspended      bool   `dynamodbav:"is_suspended" json:"isSuspended"`
	IsDeactivated    bool   `dynamodbav:"is_deactivated" json:"isDeactivated"`
	DateJoined       string `dynamodbav:"date" json:"date"`
}

type UserPersonalInfo struct {
	BirthdateChangeCount int    `dynamodbav:"birthdate_change_count" json:"birthdateChangeCount"`
	LastNicknameChange   string `dynamodbav:"last_nickname_change" json:"lastNicknameChange"`
	LastEmailChange      string `dynamodbav:"last_email_change" json:"lastEmailChange"`
	LastLogin            string `dynamodbav:"last_login" json:"lastLogin"`
	LastNameChange       string `dynamodbav:"last_name_change" json:"lastNameChange"`
	ForceChangeNickname  bool   `dynamodbav:"force_change_nickname" json:"forceChangeNickname"`
	SuspensionReason     string `dynamodbav:"suspension_reason" json:"suspensionReason"`
}

const (
	UserPkPrefix = "USER#"
	UserSkPrefix = "PROFILE"
)

func NewUser(userid string, nickname string, name string, isSuspended bool) *User {
	defaultColors := utils.GenerateRandomColorPair()

	return &User{
		UserDisplayInfo: UserDisplayInfo{
			Userid:                  userid,
			Nickname:                nickname,
			Name:                    name,
			DpUrl:                   "",
			DefaultProfilePicColors: defaultColors.Foreground + utils.AddPrefix("#", defaultColors.Background),
		},
		UserAccountInfo: UserAccountInfo{
			Bio:           "",
			IsSuspended:   isSuspended,
			IsDeactivated: false,
			DateJoined:    utils.GetTimeNow(),
		},
		UserPersonalInfo: UserPersonalInfo{
			BirthdateChangeCount: 0,
			LastNicknameChange:   "",
			LastEmailChange:      "",
			LastLogin:            utils.GetTimeNow(),
			LastNameChange:       "",
			ForceChangeNickname:  false,
			SuspensionReason:     "",
		},
		DbDescription: UserSkPrefix,
	}
}

func GetUserDisplayInfoNoId(u *User) *UserDisplayInfoNoId {
	return &UserDisplayInfoNoId{
		Nickname:                u.Nickname,
		Name:                    u.Name,
		DpUrl:                   u.DpUrl,
		DefaultProfilePicColors: u.DefaultProfilePicColors,
	}
}

func UserKey(userid string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: utils.AddPrefix(UserPkPrefix, userid)},
		"sk": &types.AttributeValueMemberS{Value: UserSkPrefix},
	}
}

func (u User) DatabaseFormat() *map[string]types.AttributeValue {
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
