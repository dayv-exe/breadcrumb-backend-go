package helpers

import (
	"breadcrumb-backend-go/utils"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func TestUserCognitoDetails(t *testing.T) {
	bDate := utils.GetTimeNow()

	expect := CognitoManagedInfo{
		Email:     "test@gmail.com",
		Birthdate: bDate,
	}

	query := func() (*cognitoidentityprovider.AdminGetUserOutput, error) {
		return &cognitoidentityprovider.AdminGetUserOutput{
			Username: aws.String("12345"),
			UserAttributes: []types.AttributeType{
				// mock cognito's response
				{Name: aws.String("email"), Value: aws.String("test@gmail.com")},
				{Name: aws.String("birthdate"), Value: aws.String(bDate)},
			},
		}, nil
	}

	result, err := getManagedInfoCustomFilter(query)

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
