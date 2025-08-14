package models

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	expected := map[string]types.AttributeValue{
		"date_joined":            &types.AttributeValueMemberS{Value: "2025-07-25T12:00:00Z"},
		"birthdate_change_count": &types.AttributeValueMemberN{Value: "2"},
		"last_nickname_change":   &types.AttributeValueMemberS{Value: "2025-07-20T10:00:00Z"},
		"last_email_change":      &types.AttributeValueMemberS{Value: "2025-07-21T08:00:00Z"},
		"last_login":             &types.AttributeValueMemberS{Value: "2025-07-25T11:00:00Z"},
		"force_change_nickname":  &types.AttributeValueMemberBOOL{Value: false},
		"suspension_reason":      &types.AttributeValueMemberS{Value: "Inappropriate content"},
		"default_pic_fg":         &types.AttributeValueMemberS{Value: userLogs.defaultProfilePicFgColor},
		"default_pic_bg":         &types.AttributeValueMemberS{Value: userLogs.defaultProfilePicBgColor},
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d keys, got %d", len(expected), len(result))
	}

	for key, expVal := range expected {
		val, exists := result[key]
		if !exists {
			t.Errorf("Missing key: %v", key)
			continue
		}

		if reflect.TypeOf(val) != reflect.TypeOf(expVal) {
			t.Errorf("For key %v: expected type: %v, but got type: %v", key, reflect.TypeOf(expVal), reflect.TypeOf(val))
			continue
		}

		if !reflect.DeepEqual(expVal, val) {
			t.Errorf("For key %v: expected %v, got %v", key, expVal, val)
			continue
		}
	}
}
