package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/omarahm3/sls/pkg/user"
)

type (
	Request  = events.APIGatewayProxyRequest
	Response = events.APIGatewayProxyResponse
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

type Handler struct {
	req       Request
	tableName string
	client    *dynamodb.DynamoDB
}

func (h *Handler) GetUser() (Response, error) {
	email := h.req.QueryStringParameters["email"]

	if email != "" {
		return getSingleUser(email, h.tableName, h.client)
	}

	return getMultipleUsers(h.tableName, h.client)
}

func (h *Handler) CreateUser() (Response, error) {
	result, err := user.CreateUser(h.req, h.tableName, h.client)
	if err != nil {
		return errorResponse(err)
	}

	return JSONResponse(http.StatusCreated, result)
}

func (h *Handler) UpdateUser() (Response, error) {
	result, err := user.UpdateUser(h.req, h.tableName, h.client)
	if err != nil {
		return errorResponse(err)
	}

	return JSONResponse(http.StatusCreated, result)
}

func (h *Handler) DeleteUser() (Response, error) {
	err := user.DeleteUser(h.req, h.tableName, h.client)
	if err != nil {
		return errorResponse(err)
	}

	return JSONResponse(http.StatusCreated, "user deleted")
}

func Prepare(req Request) (*Handler, error) {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		panic(fmt.Errorf("TABLE_NAME env variable is not set"))
	}

	dynamoClient, err := getDynamoClient()
	if err != nil {
		return nil, fmt.Errorf("could not establish db connection")
	}

	return newHandler(req, tableName, dynamoClient), nil
}

func getSingleUser(email, tableName string, dynamoClient *dynamodb.DynamoDB) (Response, error) {
	u, err := user.FetchUser(email, tableName, dynamoClient)
	if err != nil {
		return JSONResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}

	if u == nil {
		return JSONResponse(http.StatusNotFound, "user was not found")
	}

	return JSONResponse(http.StatusOK, u)
}

func getMultipleUsers(tableName string, dynamoClient *dynamodb.DynamoDB) (Response, error) {
	users, err := user.FetchUsers(tableName, dynamoClient)
	if err != nil {
		return JSONResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}

	return JSONResponse(http.StatusOK, users)
}

func errorResponse(err error) (Response, error) {
	return JSONResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
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
		Region:   aws.String(region),
		Endpoint: aws.String("http://localhost:8001"),
	}), nil
}

func newHandler(req Request, tableName string, dynamoClient *dynamodb.DynamoDB) *Handler {
	return &Handler{
		req:       req,
		tableName: tableName,
		client:    dynamoClient,
	}
}
