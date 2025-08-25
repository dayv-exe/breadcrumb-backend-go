package models

// a way to standardize responses from functions triggered by apigateway event

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type ResponseBody struct {
	Message interface{} `json:"message"`
}

// buildResponse builds a standard API Gateway proxy response
func buildResponse(statusCode int, body ResponseBody) events.APIGatewayProxyResponse {
	jsonBody, _ := json.Marshal(body)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonBody),
	}
}

func InvalidRequestErrorResponse(msg string) events.APIGatewayProxyResponse {
	if msg == "" {
		msg = "Invalid request body."
	}
	return buildResponse(400, ResponseBody{msg})
}

func UnauthorizedErrorResponse(msg string) events.APIGatewayProxyResponse {
	if msg == "" {
		msg = "Unauthorized request."
	}
	return buildResponse(401, ResponseBody{msg})
}

func NotFoundResponse(msg string) events.APIGatewayProxyResponse {
	if msg == "" {
		msg = "Resource not found."
	}

	return buildResponse(404, ResponseBody{msg})
}

func ServerSideErrorResponse(msg string, error error, errSrcDescription string) events.APIGatewayProxyResponse {
	if msg == "" {
		msg = "An error has occurred on our end, try again."
	}
	log.Println("Error source description: " + errSrcDescription + " ERROR: " + error.Error())
	return buildResponse(500, ResponseBody{msg})
}

func SuccessfulRequestResponse(msg string, createdResource bool) events.APIGatewayProxyResponse {
	if msg == "" {
		msg = "Request successful"
	}

	sCode := 200
	if createdResource {
		sCode = 201
	}

	return buildResponse(sCode, ResponseBody{msg})
}

func SuccessfulGetRequestResponse(body interface{}) events.APIGatewayProxyResponse {
	return buildResponse(200, ResponseBody{
		body,
	})
}
