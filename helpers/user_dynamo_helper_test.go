package helpers

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestIsNicknameAvailableInDynamodb(t *testing.T) {
	tests := []struct {
		name          string
		queryFn       func() (*dynamodb.GetItemOutput, error)
		wantAvailable bool
		wantErr       bool
	}{
		{
			name: "nickname taken",
			queryFn: func() (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"nickname": &types.AttributeValueMemberS{Value: "taken"},
					},
				}, nil
			},
			wantAvailable: false,
			wantErr:       false,
		},
		{
			name: "nickname not taken",
			queryFn: func() (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{
					Item: nil,
				}, nil
			},
			wantAvailable: true,
			wantErr:       false,
		},
		{
			name: "query returns error",
			queryFn: func() (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("dynamo error")
			},
			wantAvailable: false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nicknameAvailableQueryRunner(tt.queryFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.wantAvailable {
				t.Errorf("got %v, want %v", got, tt.wantAvailable)
			}
		})
	}
}
