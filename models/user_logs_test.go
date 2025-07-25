package models

import (
	"testing"
)

func TestNewUserLogs_DefaultValues(t *testing.T) {
	userLogs := NewUserLogs()

	if userLogs.BirthdateChangeCount != 0 {
		t.Errorf("expected birthdate change count to be 0, got %d", userLogs.BirthdateChangeCount)
	}

	if userLogs.forceChangeNickname != false {
		t.Errorf("expected force change nickname to be false, got true")
	}

	if userLogs.DateJoined == "" || userLogs.LastLogin == "" {
		t.Errorf("expected date joined and last login to be non-empty")
	}

	if userLogs.LastNicknameChange != "" {
		t.Errorf("expected last nickname change to be empty")
	}
}

func TestUserLogs_DatabaseFormat(t *testing.T) {
	userLogs := UserLogs{
		DateJoined:           "2025-07-25T12:00:00Z",
		BirthdateChangeCount: 2,
		LastNicknameChange:   "2025-07-20T10:00:00Z",
		LastEmailChange:      "2025-07-21T08:00:00Z",
		LastLogin:            "2025-07-25T11:00:00Z",
		forceChangeNickname:  false,
		suspensionReason:     "Inappropriate content",
	}

	result := userLogs.DatabaseFormat()

	expected := map[string]string{
		"date_joined":            "2025-07-25T12:00:00Z",
		"birthdate_change_count": "2",
		"last_nickname_change":   "2025-07-20T10:00:00Z",
		"last_email_change":      "2025-07-21T08:00:00Z",
		"last_login":             "2025-07-25T11:00:00Z",
		"force_change_nickname":  "false",
		"suspension_reason":      "Inappropriate content",
	}

	for key, expectedValue := range expected {
		if result[key] != expectedValue {
			t.Errorf("expected %s to be '%s', got '%s'", key, expectedValue, result[key])
		}
	}
}
