package models

import (
	"encoding/json"
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

	attrMap, err := user.DatabaseFormat()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expect := map[string]string{
		"pk":             "USER#123",
		"sk":             "PROFILE",
		"name":           "David Arubuike",
		"nickname":       "david",
		"dpUrl":          "https://example.com/profile.jpg",
		"is_suspended":   "true",
		"is_deactivated": "false",
	}

	for key, val := range expect {
		gotAttr, ok := attrMap[key].(*types.AttributeValueMemberS)
		if !ok {
			t.Errorf("Expected key %s to be AttributeValueMemberS, got %T", key, attrMap[key])
			continue
		}
		if gotAttr.Value != val {
			t.Errorf("Expected %s to be '%s', got '%s'", key, val, gotAttr.Value)
		}
	}

	// Check if userLogs is valid JSON
	userLogsAttr, ok := attrMap["userLogs"].(*types.AttributeValueMemberS)
	if !ok {
		t.Errorf("Expected userLogs to be AttributeValueMemberS")
	}
	var userLogs map[string]interface{}
	if err := json.Unmarshal([]byte(userLogsAttr.Value), &userLogs); err != nil {
		t.Errorf("Expected userLogs to be valid JSON, got error: %v", err)
	}
}
