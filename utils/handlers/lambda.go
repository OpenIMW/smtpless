package handlers

import (
	"os"
	"net/url"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/openimw/smtpless/utils"
)

type JsonBody map[string]string

type Response struct {
	Status  int
	Headers map[string]string
	Body    JsonBody
}

func response(res Response) (events.APIGatewayProxyResponse, error) {

	body, err := json.Marshal(res.Body)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    res.Headers,
		StatusCode: res.Status,
	}, err
}

func handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	params, err := url.ParseQuery(request.Body)

	if err != nil {
		return response(Response{
			Status: 403,
			Body: JsonBody{
				"type":    "error",
				"message": "Invalid body data",
			},
		})
	}

	email := utils.Email{
		To:   "test",
		Body: "message from: "+ params.Get("email"),
	}

	smtp := utils.SmtpConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		From:     os.Getenv("EMAIL_FROM"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}

	err = utils.Send(email, smtp)

	if err != nil {
		return response(Response{
			Status: 403,
			Body: JsonBody{
				"type":    "error",
				"message": "Cannot send email",
			},
		})
	}

	return response(Response{
		Status: 200,
		Body: JsonBody{
			"type":    "success",
			"message": "Email sent successfully.",
		},
	})
}

func LambdaHandler() {
	lambda.Start(handle)
}
