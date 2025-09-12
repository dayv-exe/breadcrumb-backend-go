package models

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestNicknameSearchIndex(t *testing.T) {

	tests := []struct {
		nickname string
		fullname string
		expected []map[string]types.AttributeValue
		wantErr  bool
	}{
		{
			"john.test",
			"johnny test",
			[]map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "jo"},
					"sk": &types.AttributeValueMemberS{Value: "john.test#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "jo"},
					"sk": &types.AttributeValueMemberS{Value: "john#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "te"},
					"sk": &types.AttributeValueMemberS{Value: "test#123"},
				},
			},
			false,
		},
		{
			"john_test",
			"johnny test",
			[]map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "jo"},
					"sk": &types.AttributeValueMemberS{Value: "john_test#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "jo"},
					"sk": &types.AttributeValueMemberS{Value: "john#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "te"},
					"sk": &types.AttributeValueMemberS{Value: "test#123"},
				},
			},
			false,
		},
		{
			"john",
			"johnny test",
			[]map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "jo"},
					"sk": &types.AttributeValueMemberS{Value: "john#123"},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nickname, func(t *testing.T) {
			us := UserSearch{
				UserId:   "123",
				Nickname: tt.nickname,
				Name:     tt.fullname,
				DpUrl:    "",
			}
			res, err := us.buildNicknameIndexItem()

			if err != nil && !tt.wantErr {
				t.Fatalf("An unexpected error has occurred: %v", err)
			} else if err == nil && tt.wantErr {
				t.Fatalf("An expected error has NOT occurred.")
			}

			if len(res) != len(tt.expected) {
				t.Fatalf("Expected %d items, got %d", len(tt.expected), len(res))
			}

			for i, expectedItem := range tt.expected {
				gotItem := res[i]

				for k, expectedVal := range expectedItem {
					gotVal, exists := gotItem[k]
					if !exists {
						t.Errorf("Item %d missing key %q", i, k)
						continue
					}

					gotStr, ok := gotVal.(*types.AttributeValueMemberS)
					if !ok {
						t.Errorf("Item %d key %q: expected AttributeValueMemberS, got %T", i, k, gotVal)
						continue
					}

					expStr := expectedVal.(*types.AttributeValueMemberS)
					if gotStr.Value != expStr.Value {
						t.Errorf("Item %d key %q: expected value %q, got %q", i, k, expStr.Value, gotStr.Value)
					}
				}
			}
		})
	}
}

func TestFullnameSearchIndex(t *testing.T) {

	tests := []struct {
		fullname string
		nickname string
		expected []map[string]types.AttributeValue
		wantErr  bool
	}{
		{
			"john test",
			"john.test",
			[]map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "jo"},
					"sk": &types.AttributeValueMemberS{Value: "john test#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "jo"},
					"sk": &types.AttributeValueMemberS{Value: "john#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#" + "te"},
					"sk": &types.AttributeValueMemberS{Value: "test#123"},
				},
			},
			false,
		},
		{
			"j😍ohn",
			"john.test",
			[]map[string]types.AttributeValue{
				{
					"pk":   &types.AttributeValueMemberS{Value: "USER_INDEX#" + "jo"},
					"sk":   &types.AttributeValueMemberS{Value: "john#123"},
					"name": &types.AttributeValueMemberS{Value: "j😍ohn"},
				},
			},
			false,
		},
		{
			"j😍 lover",
			"john.test",
			[]map[string]types.AttributeValue{
				{
					"pk":   &types.AttributeValueMemberS{Value: "USER_INDEX#" + "lo"},
					"sk":   &types.AttributeValueMemberS{Value: "lover#123"},
					"name": &types.AttributeValueMemberS{Value: "j😍 lover"},
				},
			},
			false,
		},
		{
			"😍",
			"john.test",
			[]map[string]types.AttributeValue{},
			false,
		},
		{
			"chris ronaldo yes😀",
			"john.test",
			[]map[string]types.AttributeValue{
				{
					"pk":   &types.AttributeValueMemberS{Value: "USER_INDEX#" + "ch"},
					"sk":   &types.AttributeValueMemberS{Value: "chris ronaldo yes#123"},
					"name": &types.AttributeValueMemberS{Value: "chris ronaldo yes😀"},
				},
				{
					"pk":   &types.AttributeValueMemberS{Value: "USER_INDEX#" + "ch"},
					"sk":   &types.AttributeValueMemberS{Value: "chris#123"},
					"name": &types.AttributeValueMemberS{Value: "chris ronaldo yes😀"},
				},
				{
					"pk":   &types.AttributeValueMemberS{Value: "USER_INDEX#" + "ro"},
					"sk":   &types.AttributeValueMemberS{Value: "ronaldo#123"},
					"name": &types.AttributeValueMemberS{Value: "chris ronaldo yes😀"},
				},
				{
					"pk":   &types.AttributeValueMemberS{Value: "USER_INDEX#" + "ye"},
					"sk":   &types.AttributeValueMemberS{Value: "yes#123"},
					"name": &types.AttributeValueMemberS{Value: "chris ronaldo yes😀"},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.fullname, func(t *testing.T) {
			us := UserSearch{
				UserId:   "123",
				Name:     tt.fullname,
				Nickname: tt.nickname,
				DpUrl:    "",
			}
			res, err := us.buildFullnameIndexItem()

			if err != nil && !tt.wantErr {
				t.Fatalf("An unexpected error has occurred: %v", err)
			} else if err == nil && tt.wantErr {
				t.Fatalf("An expected error has NOT occurred.")
			}

			if len(res) != len(tt.expected) {
				t.Fatalf("Expected %d items, got %d", len(tt.expected), len(res))
			}

			for i, expectedItem := range tt.expected {
				gotItem := res[i]

				for k, expectedVal := range expectedItem {
					gotVal, exists := gotItem[k]
					if !exists {
						t.Errorf("Item %d missing key %q", i, k)
						continue
					}

					gotStr, ok := gotVal.(*types.AttributeValueMemberS)
					if !ok {
						t.Errorf("Item %d key %q: expected AttributeValueMemberS, got %T", i, k, gotVal)
						continue
					}

					expStr := expectedVal.(*types.AttributeValueMemberS)
					if gotStr.Value != expStr.Value {
						t.Errorf("Item %d key %q: expected value %q, got %q", i, k, expStr.Value, gotStr.Value)
					}
				}
			}
		})
	}
}

func TestGetUserSearchIndexesKeysNickName(t *testing.T) {
	tests := []struct {
		userid       string
		name         string
		nickname     string
		expectedKeys []map[string]types.AttributeValue
	}{
		{
			"123",
			"john",
			"john",
			[]map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jo"},
					"sk": &types.AttributeValueMemberS{Value: "john#123"},
				},
			},
		},
		{
			"123",
			"j😍lover",
			"john_love",
			[]map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jo"},
					"sk": &types.AttributeValueMemberS{Value: "john_love#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jo"},
					"sk": &types.AttributeValueMemberS{Value: "john#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#lo"},
					"sk": &types.AttributeValueMemberS{Value: "love#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jl"},
					"sk": &types.AttributeValueMemberS{Value: "jlover#123"},
				},
			},
		},
		{
			"123",
			"chris ronaldo yes😀",
			"cry",
			[]map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#cr"},
					"sk": &types.AttributeValueMemberS{Value: "cry#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#ch"},
					"sk": &types.AttributeValueMemberS{Value: "chris ronaldo yes#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#ch"},
					"sk": &types.AttributeValueMemberS{Value: "chris#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#ro"},
					"sk": &types.AttributeValueMemberS{Value: "ronaldo#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#ye"},
					"sk": &types.AttributeValueMemberS{Value: "yes#123"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nickname, func(t *testing.T) {
			u := UserSearch{
				UserId:   "123",
				Nickname: tt.nickname,
				Name:     tt.name,
			}
			indexes, _ := u.BuildSearchIndexes()
			result := GetUserSearchIndexesKeys(indexes)

			if len(result) != len(tt.expectedKeys) {
				t.Fatalf("Expected %d items, got %d", len(tt.expectedKeys), len(result))
			}

			for i, expectedItem := range tt.expectedKeys {
				gotItem := result[i]

				for k, expectedVal := range expectedItem {
					gotVal, exists := gotItem[k]
					if !exists {
						t.Errorf("Item %d missing key %q", i, k)
						continue
					}

					gotStr, ok := gotVal.(*types.AttributeValueMemberS)
					if !ok {
						t.Errorf("Item %d key %q: expected AttributeValueMemberS, got %T", i, k, gotVal)
						continue
					}

					expStr := expectedVal.(*types.AttributeValueMemberS)
					if gotStr.Value != expStr.Value {
						t.Errorf("Item %d key %q: expected value %q, got %q", i, k, expStr.Value, gotStr.Value)
					}
				}
			}
		})
	}
}

func TestDeconstructSearchIndexItem(t *testing.T) {
	tests := []struct {
		expectedName     string
		expectedNickname string
		searchDbItem     map[string]types.AttributeValue
		wantErr          bool
	}{
		{
			"johnny",
			"johnny.test",
			map[string]types.AttributeValue{
				"pk":       &types.AttributeValueMemberS{Value: "NICKNAME#jo"},
				"sk":       &types.AttributeValueMemberS{Value: "NICKNAME#johnny.test"},
				"name":     &types.AttributeValueMemberS{Value: "johnny"},
				"nickname": &types.AttributeValueMemberS{Value: "johnny.test"},
			},
			false,
		},
		{
			"j🤣ohnny",
			"johnny.test",
			map[string]types.AttributeValue{
				"pk":       &types.AttributeValueMemberS{Value: "NICKNAME#jo"},
				"sk":       &types.AttributeValueMemberS{Value: "NICKNAME#johnny.test"},
				"name":     &types.AttributeValueMemberS{Value: "j🤣ohnny"},
				"nickname": &types.AttributeValueMemberS{Value: "johnny.test"},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.expectedName, func(t *testing.T) {
			res, err := DeconstructSearchIndexItem(&tt.searchDbItem)

			if err != nil && !tt.wantErr {
				t.Fatalf("An unexpected error has occurred: %v", err)
			} else if err == nil && tt.wantErr {
				t.Fatalf("An expected error has NOT occurred.")
			}

			if res.Name != tt.expectedName {
				t.Fatalf("Expected name to be %v got %v instead", tt.expectedName, res.Name)
			}

			if res.Nickname != tt.expectedNickname {
				t.Fatalf("Expected nickname to be %v got %v instead", tt.expectedNickname, res.Nickname)
			}
		})
	}
}
