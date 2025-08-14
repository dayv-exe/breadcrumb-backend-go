package models

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestUser_DatabaseFormat(t *testing.T) {
	user := NewUser(
		"123",
		"David",
		"David Arubuike",
		true,
	)

	user.DpUrl = "https://example.com/profile.jpg"

	result := user.DatabaseFormat()

	expected := map[string]types.AttributeValue{
		"pk":             &types.AttributeValueMemberS{Value: "USER#123"},
		"sk":             &types.AttributeValueMemberS{Value: "PROFILE"},
		"name":           &types.AttributeValueMemberS{Value: "David Arubuike"},
		"nickname":       &types.AttributeValueMemberS{Value: "david"},
		"bio":            &types.AttributeValueMemberS{Value: ""},
		"dpUrl":          &types.AttributeValueMemberS{Value: "https://example.com/profile.jpg"},
		"is_suspended":   &types.AttributeValueMemberBOOL{Value: true},
		"is_deactivated": &types.AttributeValueMemberBOOL{Value: false},
		"user_logs":      &types.AttributeValueMemberM{Value: user.UserLogs.DatabaseFormat()},
	}

	if len(expected) != len(result) {
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

		if !reflect.DeepEqual(val, expVal) {
			t.Errorf("For key %v: expected %v, got %v", key, expVal, val)
			continue
		}
	}
}
