package server_mocks

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Mocks struct{}

func New() *Mocks {
	return &Mocks{}
}

func (*Mocks) Lambda() func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		response := events.APIGatewayProxyResponse{}
		response.Body = `{"message":"hello"}`
		response.StatusCode = http.StatusOK
		return response, nil
	}
}
