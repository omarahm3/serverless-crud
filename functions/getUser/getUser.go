package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/omarahm3/sls/pkg/handlers"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		panic(fmt.Errorf("TABLE_NAME env variable is not set"))
	}

	dynamoClient, err := getDynamoClient()
	if err != nil {
		return handlers.JSONResponse(http.StatusInternalServerError, "could not establish db connection")
	}

	return handlers.GetUser(req, tableName, dynamoClient)
}

func main() {
	lambda.Start(handler)
}

func getDynamoClient() (*dynamodb.DynamoDB, error) {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return dynamodb.New(awsSession, &aws.Config{
		Region:   aws.String("eu-central-1"),
		Endpoint: aws.String("http://localhost:8001"),
	}), nil
}
