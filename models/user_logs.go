package models

import (
	"breadcrumb-backend-go/utils"
	"fmt"
)

type UserLogs struct {
	DateJoined           string
	BirthdateChangeCount int
	LastNicknameChange   string
	LastEmailChange      string
	LastLogin            string
}

func NewUserLogs() UserLogs {
	return UserLogs{
		DateJoined:           utils.GetTimeNow(),
		BirthdateChangeCount: 0,
		LastLogin:            utils.GetTimeNow(),
	}
}

func (ul UserLogs) DatabaseFormat() map[string]string {
	return map[string]string{
		"date_joined":            ul.DateJoined,
		"birthdate_change_count": fmt.Sprintf("%v", ul.BirthdateChangeCount),
		"last_nickname_change":   ul.LastNicknameChange,
		"last_email_change":      ul.LastEmailChange,
		"last_login":             ul.LastLogin,
	}
}
