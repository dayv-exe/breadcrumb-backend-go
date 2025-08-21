package models

import (
	"breadcrumb-backend-go/utils"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	dbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func TestUser_DatabaseFormat(t *testing.T) {
	user := NewUser(
		"123",
		"david",
		"David Arubuike",
		true,
	)

	user.DpUrl = "https://example.com/profile.jpg"

	result := user.DatabaseFormat()

	result["user_logs"] = nil

	expected := map[string]dbTypes.AttributeValue{
		"pk":             &dbTypes.AttributeValueMemberS{Value: "USER#123"},
		"sk":             &dbTypes.AttributeValueMemberS{Value: "PROFILE"},
		"name":           &dbTypes.AttributeValueMemberS{Value: "David Arubuike"},
		"nickname":       &dbTypes.AttributeValueMemberS{Value: "david"},
		"bio":            &dbTypes.AttributeValueMemberS{Value: ""},
		"dpUrl":          &dbTypes.AttributeValueMemberS{Value: "https://example.com/profile.jpg"},
		"is_suspended":   &dbTypes.AttributeValueMemberBOOL{Value: true},
		"is_deactivated": &dbTypes.AttributeValueMemberBOOL{Value: false},
		"user_logs":      nil,
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

func TestUserCognitoDetails(t *testing.T) {
	bDate := utils.GetTimeNow()

	expect := CognitoInfo{
		Email:     "test@gmail.com",
		Birthdate: bDate,
	}

	query := func() (*cognitoidentityprovider.AdminGetUserOutput, error) {
		return &cognitoidentityprovider.AdminGetUserOutput{
			Username: aws.String("12345"),
			UserAttributes: []cognitoTypes.AttributeType{
				// mock cognito's response
				{Name: aws.String("email"), Value: aws.String("test@gmail.com")},
				{Name: aws.String("birthdate"), Value: aws.String(bDate)},
			},
		}, nil
	}

	result, err := getCognitoInfoQueryRunner(query)

	if err != nil {
		t.Fatalf("An unexpected error has occurred: %v", err)
	}

	if result.Email != expect.Email {
		t.Fatalf("For key email: expected %v, got %v", expect.Email, result.Email)
	}

	if result.Birthdate != expect.Birthdate {
		t.Fatalf("For key birthdate: expected %v, got %v", expect.Birthdate, result.Birthdate)
	}
}
