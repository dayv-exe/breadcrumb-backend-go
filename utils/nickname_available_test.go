package utils

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestIsNicknameTakenInDynamodb(t *testing.T) {
	tests := []struct {
		name      string
		queryFn   func() (*dynamodb.QueryOutput, error)
		wantTaken bool
		wantErr   bool
	}{
		{
			name: "nickname taken",
			queryFn: func() (*dynamodb.QueryOutput, error) {
				return &dynamodb.QueryOutput{
					Items: []map[string]ddbTypes.AttributeValue{
						{"nickname": &ddbTypes.AttributeValueMemberS{Value: "taken"}},
					},
				}, nil
			},
			wantTaken: true,
			wantErr:   false,
		},
		{
			name: "nickname not taken",
			queryFn: func() (*dynamodb.QueryOutput, error) {
				return &dynamodb.QueryOutput{
					Items: []map[string]ddbTypes.AttributeValue{},
				}, nil
			},
			wantTaken: false,
			wantErr:   false,
		},
		{
			name: "query returns error",
			queryFn: func() (*dynamodb.QueryOutput, error) {
				return nil, errors.New("dynamo error")
			},
			wantTaken: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NicknameAvailabilityCheck(tt.queryFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.wantTaken {
				t.Errorf("got %v, want %v", got, tt.wantTaken)
			}
		})
	}
}
