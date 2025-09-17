package models

import (
	"breadcrumb-backend-go/utils"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDatabaseFormat(t *testing.T) {
	d := utils.GetTimeNow()

	expected := map[string]types.AttributeValue{
		"pk":   &types.AttributeValueMemberS{Value: "USER#123"},
		"sk":   &types.AttributeValueMemberS{Value: "FRIEND#321"},
		"date": &types.AttributeValueMemberS{Value: d},
	}

	f := friend{
		User1Id: "123",
		User2Id: "321",
		Date:    d,
	}
	res, _ := f.DatabaseFormat()
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
