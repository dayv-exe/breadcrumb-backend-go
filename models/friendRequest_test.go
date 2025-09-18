package models

import (
	"breadcrumb-backend-go/utils"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestFriendRequestDbFormat(t *testing.T) {
	d := utils.GetTimeNow()
	expected := map[string]types.AttributeValue{
		"pk":             &types.AttributeValueMemberS{Value: "USER#rec"},
		"sk":             &types.AttributeValueMemberS{Value: "FRIEND_REQUEST_FROM#send"},
		"date":           &types.AttributeValueMemberS{Value: d},
		"name":           &types.AttributeValueMemberS{Value: "test"},
		"nickname":       &types.AttributeValueMemberS{Value: "test"},
		"dp_url":         &types.AttributeValueMemberS{Value: ""},
		"default_pic_fg": &types.AttributeValueMemberS{Value: ""},
		"default_pic_bg": &types.AttributeValueMemberS{Value: ""},
	}

	fr := FriendRequest{
		RecipientId:     "rec",
		SenderId:        "send",
		Date:            d,
		SendersName:     "test",
		SendersNickname: "test",
		SendersDpUrl:    "",
		SendersFgCol:    "",
		SendersBgCol:    "",
	}
	res, _ := fr.DatabaseFormat()
	result := *res

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

		if !reflect.DeepEqual(val, expVal) {
			t.Errorf("For key %v: expected %v, got %v", key, expVal, val)
			continue
		}
	}
}

func TestConvertToFriendRequest(t *testing.T) {
	d := utils.GetTimeNow()
	expected := FriendRequest{
		RecipientId: "rec",
		SenderId:    "send",
		Date:        d,
	}

	item := map[string]types.AttributeValue{
		"pk":   &types.AttributeValueMemberS{Value: "USER#rec"},
		"sk":   &types.AttributeValueMemberS{Value: "FRIEND_REQUEST_FROM#send"},
		"date": &types.AttributeValueMemberS{Value: d},
	}

	result, _ := ConvertToFriendRequest(&item)

	if !reflect.DeepEqual(*result, expected) {
		t.Errorf("Result: %v does not match expected: %v", result, expected)
	}
}
