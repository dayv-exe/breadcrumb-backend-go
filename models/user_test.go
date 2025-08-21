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

var testUserDynamo = map[string]dbTypes.AttributeValue{
	"pk":                     &dbTypes.AttributeValueMemberS{Value: "USER#123"},
	"sk":                     &dbTypes.AttributeValueMemberS{Value: "PROFILE"},
	"nickname":               &dbTypes.AttributeValueMemberS{Value: "test"},
	"name":                   &dbTypes.AttributeValueMemberS{Value: "test test"},
	"bio":                    &dbTypes.AttributeValueMemberS{Value: ""},
	"dpUrl":                  &dbTypes.AttributeValueMemberS{Value: ""},
	"is_suspended":           &dbTypes.AttributeValueMemberBOOL{Value: false},
	"is_deactivated":         &dbTypes.AttributeValueMemberBOOL{Value: false},
	"date_joined":            &dbTypes.AttributeValueMemberS{Value: utils.GetTimeNow()},
	"birthdate_change_count": &dbTypes.AttributeValueMemberN{Value: "0"},
	"last_nickname_change":   &dbTypes.AttributeValueMemberS{Value: ""},
	"last_email_change":      &dbTypes.AttributeValueMemberS{Value: ""},
	"last_login":             &dbTypes.AttributeValueMemberS{Value: utils.GetTimeNow()},
	"force_change_nickname":  &dbTypes.AttributeValueMemberBOOL{Value: false},
	"suspension_reason":      &dbTypes.AttributeValueMemberS{Value: ""},
	"default_pic_fg":         &dbTypes.AttributeValueMemberS{Value: ""},
	"default_pic_bg":         &dbTypes.AttributeValueMemberS{Value: ""},
}

func TestUser_DatabaseFormat(t *testing.T) {
	user := NewUser(
		"123",
		"test",
		"test test",
		false,
	)

	result := user.DatabaseFormat()

	result["default_pic_fg"] = &dbTypes.AttributeValueMemberS{Value: ""}
	result["default_pic_bg"] = &dbTypes.AttributeValueMemberS{Value: ""}

	if len(testUserDynamo) != len(result) {
		t.Fatalf("Expected %d keys, got %d", len(testUserDynamo), len(result))
	}

	for key, expVal := range testUserDynamo {
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

func TestConvertToUser(t *testing.T) {
	expect := User{
		Userid:        "123",
		Nickname:      "test",
		Name:          "test",
		Bio:           "",
		DpUrl:         "",
		IsSuspended:   false,
		IsDeactivated: false,
	}

	result, err := convertToUser(map[string]dbTypes.AttributeValue{
		"pk":             &dbTypes.AttributeValueMemberS{Value: "USER#123"},
		"nickname":       &dbTypes.AttributeValueMemberS{Value: "test"},
		"name":           &dbTypes.AttributeValueMemberS{Value: "test"},
		"bio":            &dbTypes.AttributeValueMemberS{Value: ""},
		"dpUrl":          &dbTypes.AttributeValueMemberS{Value: ""},
		"is_suspended":   &dbTypes.AttributeValueMemberBOOL{Value: false},
		"is_deactivated": &dbTypes.AttributeValueMemberBOOL{Value: false},
	})

	if err != nil {
		t.Fatalf("An unexpected error occurred %v", err.Error())
	}

	if !reflect.DeepEqual(*result, expect) {
		t.Errorf("Result: %v does not match expected: %v", result, expect)
	}
}
