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
		"David",
		"David Arubuike",
		true,
	)

	user.DpUrl = "https://example.com/profile.jpg"

	result := user.DatabaseFormat()

	expected := map[string]dbTypes.AttributeValue{
		"pk":             &dbTypes.AttributeValueMemberS{Value: "USER#123"},
		"sk":             &dbTypes.AttributeValueMemberS{Value: "PROFILE"},
		"name":           &dbTypes.AttributeValueMemberS{Value: "David Arubuike"},
		"nickname":       &dbTypes.AttributeValueMemberS{Value: "david"},
		"bio":            &dbTypes.AttributeValueMemberS{Value: ""},
		"dpUrl":          &dbTypes.AttributeValueMemberS{Value: "https://example.com/profile.jpg"},
		"is_suspended":   &dbTypes.AttributeValueMemberBOOL{Value: true},
		"is_deactivated": &dbTypes.AttributeValueMemberBOOL{Value: false},
		"user_logs":      &dbTypes.AttributeValueMemberM{Value: user.UserLogs.DatabaseFormat()},
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
		Cognito: map[string]string{
			"email":     "test@gmail.com",
			"birthdate": bDate,
		},
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

	if len(result.Cognito) != len(expect.Cognito) {
		t.Fatalf("Expected %d keys, got %d", len(expect.Cognito), len(result.Cognito))
	}

	for key, expVal := range expect.Cognito {
		val, exists := result.Cognito[key]

		// check for missing keys
		if !exists {
			t.Fatalf("Missing key: %v", key)
			continue
		}

		// check that values match
		if !reflect.DeepEqual(val, expVal) {
			t.Errorf("For key %v: expected %v, got %v", key, expVal, val)
			continue
		}

		// check that the data types match
		if reflect.TypeOf(val) != reflect.TypeOf(expVal) {
			t.Errorf("For key %v: expected type: %v, but got type: %v", key, reflect.TypeOf(expVal), reflect.TypeOf(val))
			continue
		}
	}
}
