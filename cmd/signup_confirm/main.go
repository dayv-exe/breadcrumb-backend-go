package main

import (
	"breadcrumb-backend-go/handlers/auth"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(auth.HandlePostConfirmation)
}
