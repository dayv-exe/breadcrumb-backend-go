package utils

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
					Item: map[string]ddbTypes.AttributeValue{
						"nickname": &ddbTypes.AttributeValueMemberS{Value: "taken"},
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

func TestIsNicknameValid(t *testing.T) {
	tests := []struct {
		nickname string
		valid    bool
	}{
		{"john_doe", true},
		{"_john", false},
		{"john..doe", false},
		{"j", false},
		{"averylongnamethatshouldfail", false},
		{"john.doe_", false},
		{"john__doe", false},
		{"j.doe", true},
		{"j_doe", true},
		{"JohnDoe99", true},
		{"john_.doe", false},
		{"john.doe_", false},
		{".johndoe", false},
		{"4ksf_sqmd1", true},
		{"john.doe001234ed", false},
		{"14792384913", true},
	}

	for _, tt := range tests {
		t.Run(tt.nickname, func(t *testing.T) {
			got := NicknameValid(tt.nickname)
			if got != tt.valid {
				t.Errorf("isNicknameValid(%q) = %v; want %v", tt.nickname, got, tt.valid)
			}
		})
	}
}
