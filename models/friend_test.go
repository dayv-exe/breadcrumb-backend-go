package models

import (
	"breadcrumb-backend-go/utils"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestFriendDatabaseFormat(t *testing.T) {
	d := utils.GetTimeNow()

	expected := map[string]types.AttributeValue{
		"pk":                 &types.AttributeValueMemberS{Value: "USER#123"},
		"sk":                 &types.AttributeValueMemberS{Value: "FRIEND#321"},
		"name":               &types.AttributeValueMemberS{Value: "sender"},
		"nickname":           &types.AttributeValueMemberS{Value: "sndr"},
		"dpUrl":              &types.AttributeValueMemberS{Value: ""},
		"default_pic_colors": &types.AttributeValueMemberS{Value: ""},
		"date":               &types.AttributeValueMemberS{Value: d},
	}

	otherUser := NewUser("321", "sndr", "sender", false)
	otherUser.DefaultProfilePicColors = ""

	result, _ := NewFriendship("123", otherUser).DatabaseFormat()

	if len(*result) != len(expected) {
		t.Fatalf("Expected %d keys, got %d", len(expected), len(*result))
	}

	for key, expVal := range expected {
		val, exists := (*result)[key]
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

func TestFriendToUserInfoStruct(t *testing.T) {
	friendDbItem := []map[string]types.AttributeValue{
		{
			"pk":                 &types.AttributeValueMemberS{Value: "123"},
			"sk":                 &types.AttributeValueMemberS{Value: "321"},
			"nickname":           &types.AttributeValueMemberS{Value: "other"},
			"name":               &types.AttributeValueMemberS{Value: "other"},
			"dpUrl":              &types.AttributeValueMemberS{Value: ""},
			"default_pic_colors": &types.AttributeValueMemberS{Value: ""},
			"date":               &types.AttributeValueMemberS{Value: ""},
		},
	}

	expect := UserDisplayInfo{
		Userid:                  "321",
		Nickname:                "other",
		Name:                    "other",
		DpUrl:                   "",
		DefaultProfilePicColors: "",
	}

	result, _ := FriendItemsToUserDisplayStructs(&friendDbItem)

	if !reflect.DeepEqual((*result)[0], expect) {
		t.Errorf("expected %v, got %v", expect, (*result)[0])
	}
}
