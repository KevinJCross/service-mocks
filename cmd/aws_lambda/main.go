package main

import (
	"github.com/KevinJCross/server_mocks"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(server_mocks.New().Lambda())
}
