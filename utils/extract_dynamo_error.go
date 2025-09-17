package utils

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PrintTransactWriteCancellationReason(err error) {
	// Check for transaction cancellation reasons
	var tce *types.TransactionCanceledException
	if errors.As(err, &tce) {
		for i, reason := range tce.CancellationReasons {
			fmt.Printf("add user Cancellation %d: Code=%s, Message=%s\n",
				i,
				aws.ToString(reason.Code),
				aws.ToString(reason.Message),
			)
		}
	}
}
