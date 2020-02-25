package handlers

import (
	// "github.com/openimw/smtpless/utils"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
)


func handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		Body: "Hello",
		StatusCode: 200,
	}, nil
}

func LambdaHandler() {
	lambda.Start(handle)
}

