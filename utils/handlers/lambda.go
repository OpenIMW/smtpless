package handlers

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/openimw/smtpless/utils"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
	"net/url"
	"os"
	"time"
)

type JsonBody struct {
	Success bool   `json:"success"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Status  int
	Headers map[string]string
	Body    JsonBody
}

func respond(res Response) (events.APIGatewayProxyResponse, error) {

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
		return respond(Response{
			Status: 403,
			Body: JsonBody{
				Type:    "data",
				Success: false,
				Message: "Invalid body data",
			},
		})
	}

	var captcha recaptcha.ReCAPTCHA

	if captcha, err = recaptcha.NewReCAPTCHA(os.Getenv("RECAPTCHA_SECRET"), recaptcha.V3, 10*time.Second); err == nil {

		err = captcha.Verify(params.Get("recaptcha_response"))
	}

	if err != nil {
		return respond(Response{
			Status: 403,
			Body: JsonBody{
				Type:    "recaptcha",
				Success: false,
				Message: err.Error(),
			},
		})
	}

	email := utils.Email{
		To:   "test",
		Body: "message from: " + params.Get("email"),
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
		return respond(Response{
			Status: 403,
			Body: JsonBody{
				Type:    "mail",
				Success: false,
				Message: "Cannot send email",
			},
		})
	}

	return respond(Response{
		Status: 200,
		Body: JsonBody{
			Type:    "mail",
			Success: true,
			Message: "Email sent successfully.",
		},
	})
}

func LambdaHandler() {
	lambda.Start(handle)
}
