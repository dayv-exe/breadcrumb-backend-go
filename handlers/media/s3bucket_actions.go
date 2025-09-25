package media

import (
	"breadcrumb-backend-go/constants"
	"breadcrumb-backend-go/helpers"
	"breadcrumb-backend-go/models"
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3BucketActionsDependencies struct {
	S3Client   *s3.Client
	BucketName string
}

func handleGeneratePresignedUrl(req *events.APIGatewayProxyRequest, s3Client *s3.Client, bucketName string, ctx context.Context) (events.APIGatewayProxyResponse, error) {
	var input helpers.PresignRequest
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		log.Println("unable to convert req body to request struct")
		return models.InvalidRequestErrorResponse(""), nil
	}
	if input.Key == "" {
		log.Println("key is empty")
		return models.InvalidRequestErrorResponse(""), nil
	}

	// to check file type
	// more checks and validations will be done after upload completes with s3 event lambda
	ext := strings.ToLower(filepath.Ext(input.Key))
	if _, ok := constants.ALLOWED_FILE_TYPES[ext]; !ok {
		log.Println("INvalid file type provided: " + ext)
		return models.InvalidRequestErrorResponse("incorrect file type."), nil
	}

	if !strings.HasPrefix(input.Key, "uploads/") {
		log.Println("trying to upload to path other than uploads " + input.Key)
		return models.InvalidRequestErrorResponse(""), nil
	}

	mediaHelper := helpers.MediaHelper{
		S3Client:   s3Client,
		BucketName: bucketName,
		Ctx:        ctx,
	}

	result, err := mediaHelper.GeneratePresignedUrl(&input)
	if err != nil {
		return models.ServerSideErrorResponse("failed to generate presigned url", err, "error while generating presigned url"), nil
	}

	return models.SuccessfulGetRequestResponse(result), nil
}

func (deps *S3BucketActionsDependencies) HandleStorageActions(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	action, ok := req.PathParameters["action"]
	if !ok || action == "" {
		return models.InvalidRequestErrorResponse(""), nil
	}

	switch action {
	case "presign-url":
		return handleGeneratePresignedUrl(req, deps.S3Client, deps.BucketName, ctx)

	default:
		return models.InvalidRequestErrorResponse("invalid action requested"), nil
	}

}
