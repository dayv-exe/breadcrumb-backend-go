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
		"pk":                 &types.AttributeValueMemberS{Value: "USER#rec"},
		"sk":                 &types.AttributeValueMemberS{Value: "FRIEND_REQUEST_FROM#send"},
		"date":               &types.AttributeValueMemberS{Value: d},
		"name":               &types.AttributeValueMemberS{Value: "test"},
		"nickname":           &types.AttributeValueMemberS{Value: "test"},
		"dpUrl":              &types.AttributeValueMemberS{Value: ""},
		"default_pic_colors": &types.AttributeValueMemberS{Value: ""},
	}

	fr := friendRequest{
		RecipientId: "rec",
		SenderId:    "send",
		Date:        d,
		UserDisplayInfoNoId: UserDisplayInfoNoId{
			Name:                    "test",
			Nickname:                "test",
			DpUrl:                   "",
			DefaultProfilePicColors: "",
		},
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
	expected := UserDisplayInfo{
		Userid:                  "send",
		Nickname:                "send",
		Name:                    "send",
		DpUrl:                   "",
		DefaultProfilePicColors: "",
	}

	item := []map[string]types.AttributeValue{
		{
			"pk":                 &types.AttributeValueMemberS{Value: "USER#rec"},
			"sk":                 &types.AttributeValueMemberS{Value: "FRIEND_REQUEST_FROM#send"},
			"nickname":           &types.AttributeValueMemberS{Value: "send"},
			"name":               &types.AttributeValueMemberS{Value: "send"},
			"dpUrl":              &types.AttributeValueMemberS{Value: ""},
			"default_pic_colors": &types.AttributeValueMemberS{Value: ""},
			"date":               &types.AttributeValueMemberS{Value: d},
		},
	}

	results, _ := FriendRequestItemsToUserDisplayStructs(&item)
	if !reflect.DeepEqual((*results)[0], expected) {
		t.Errorf("Result: %v does not match expected: %v", (*results)[0], expected)
	}
}

func TestFriendRequestToUserInfoStruct(t *testing.T) {
	friendReqDbItem := []map[string]types.AttributeValue{
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

	result, _ := FriendRequestItemsToUserDisplayStructs(&friendReqDbItem)

	if !reflect.DeepEqual((*result)[0], expect) {
		t.Errorf("expected %v, got %v", expect, (*result)[0])
	}
}
