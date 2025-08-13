package models

import (
	"breadcrumb-backend-go/utils"
	"fmt"
	"strings"
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

func (ul UserLogs) DatabaseFormat() map[string]string {
	return map[string]string{
		"date_joined":            ul.DateJoined,
		"birthdate_change_count": fmt.Sprintf("%v", ul.BirthdateChangeCount),
		"last_nickname_change":   ul.LastNicknameChange,
		"last_email_change":      ul.LastEmailChange,
		"last_login":             ul.LastLogin,
		"force_change_nickname":  strings.ToLower(fmt.Sprint(ul.forceChangeNickname)),
		"default_pic_fg":         ul.defaultProfilePicFgColor,
		"default_pic_bg":         ul.defaultProfilePicBgColor,
		"suspension_reason":      ul.suspensionReason,
	}
}
