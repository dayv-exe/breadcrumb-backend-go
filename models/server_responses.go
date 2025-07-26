package models

// a way to standardize responses from functions triggered by apigateway event

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

// buildResponse builds a standard API Gateway proxy response
func buildResponse(statusCode int, body map[string]string) events.APIGatewayProxyResponse {
	jsonBody, _ := json.Marshal(body)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonBody),
	}
}

// responseMessage formats a message or error into a response
func responseMessage(statusCode int, message string) events.APIGatewayProxyResponse {
	var msgKey string
	if statusCode >= 300 {
		msgKey = "error"
	} else {
		msgKey = "message"
	}

	return buildResponse(statusCode, map[string]string{
		msgKey: message,
	})
}

func InvalidRequestErrorResponse(msg string) events.APIGatewayProxyResponse {
	if msg == "" {
		msg = "Invalid request body."
	}
	return responseMessage(400, msg)
}

func ServerSideErrorResponse(msg string, error error) events.APIGatewayProxyResponse {
	if msg == "" {
		msg = "An error has occurred on our end, try again."
	}
	log.Println(error)
	return responseMessage(500, msg)
}

func SuccessfulRequestResponse(msg string, createdResource bool) events.APIGatewayProxyResponse {
	if msg == "" {
		msg = "Request successful"
	}

	sCode := 200
	if createdResource {
		sCode = 201
	}

	return responseMessage(sCode, msg)
}
