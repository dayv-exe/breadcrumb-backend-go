package models

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestGetUserSearchIndexesKeysNickName(t *testing.T) {
	tests := []struct {
		userid       string
		name         string
		nickname     string
		expectedKeys []map[string]types.AttributeValue
	}{
		{
			userid:   "123",
			name:     "john",
			nickname: "john",
			expectedKeys: []map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jo"},
					"sk": &types.AttributeValueMemberS{Value: "john#123"},
				},
			},
		},
		{
			userid:   "123",
			name:     "jane doe",
			nickname: "jd_247",
			expectedKeys: []map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#ja"},
					"sk": &types.AttributeValueMemberS{Value: "jane#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#do"},
					"sk": &types.AttributeValueMemberS{Value: "doe#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jd"},
					"sk": &types.AttributeValueMemberS{Value: "jd#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#24"},
					"sk": &types.AttributeValueMemberS{Value: "247#123"},
				},
			},
		},
		{
			userid:   "123",
			name:     "john john",
			nickname: "john.john",
			expectedKeys: []map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jo"},
					"sk": &types.AttributeValueMemberS{Value: "john#123"},
				},
			},
		},
		{
			userid:   "123",
			name:     "j❤️t",
			nickname: "j_ohnjohn",
			expectedKeys: []map[string]types.AttributeValue{
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#oh"},
					"sk": &types.AttributeValueMemberS{Value: "ohnjohn#123"},
				},
			},
		},
		{
			userid:   "123",
			name:     "jt.vence",
			nickname: "jo_hntohn",
			expectedKeys: []map[string]types.AttributeValue{
				// names
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jt"},
					"sk": &types.AttributeValueMemberS{Value: "jt#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#ve"},
					"sk": &types.AttributeValueMemberS{Value: "vence#123"},
				},

				//nicknames
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#jo"},
					"sk": &types.AttributeValueMemberS{Value: "jo#123"},
				},
				{
					"pk": &types.AttributeValueMemberS{Value: "USER_INDEX#hn"},
					"sk": &types.AttributeValueMemberS{Value: "hntohn#123"},
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

func TestBuildUserSearchIndexes(t *testing.T) {
	tests := []struct {
		userId   string
		name     string
		nickname string
		expect   []map[string]types.AttributeValue
	}{
		{
			userId:   "123",
			name:     "johnny test",
			nickname: "j.test",
			expect: []map[string]types.AttributeValue{
				{
					"pk":       &types.AttributeValueMemberS{Value: "USER_INDEX#jo"},
					"sk":       &types.AttributeValueMemberS{Value: "johnny#123"},
					"userid":   &types.AttributeValueMemberS{Value: "123"},
					"name":     &types.AttributeValueMemberS{Value: "johnny test"},
					"nickname": &types.AttributeValueMemberS{Value: "j.test"},
					"dp_url":   &types.AttributeValueMemberS{Value: ""},
				},
				{
					"pk":       &types.AttributeValueMemberS{Value: "USER_INDEX#te"},
					"sk":       &types.AttributeValueMemberS{Value: "test#123"},
					"userid":   &types.AttributeValueMemberS{Value: "123"},
					"name":     &types.AttributeValueMemberS{Value: "johnny test"},
					"nickname": &types.AttributeValueMemberS{Value: "j.test"},
					"dp_url":   &types.AttributeValueMemberS{Value: ""},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserSearch{
				UserId:   tt.userId,
				Nickname: tt.nickname,
				Name:     tt.name,
				DpUrl:    "",
			}

			result, err := u.BuildSearchIndexes()
			if err != nil {
				t.Fatalf("An unexpected error has occurred: %v", err)
			}

			if len(result) != len(tt.expect) {
				t.Fatalf("Expected %v items but got %v instead", len(tt.expect), len(result))
			}

			for i, expectedItem := range tt.expect {
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
