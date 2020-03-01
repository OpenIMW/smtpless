package handlers

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/openimw/smtpless/utils"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
	"io/ioutil"
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

type DistConfig struct {
	Destinations []struct {
		Token string   `json:"token"`
		Host  string   `json:"host"`
		To    []string `json:"to"`
	} `json:"destinations"`
}

func respond(res Response) (events.APIGatewayProxyResponse, error) {

	body, err := json.Marshal(res.Body)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    res.Headers,
		StatusCode: res.Status,
	}, err
}

func readFile(name string) ([]byte, error) {

	cwd, err := os.Getwd()

	if err != nil {

		return []byte{}, err
	}

	return ioutil.ReadFile(cwd + name)
}

// TODO: resolve host and verify hmac hash
func resolveDist(c DistConfig) (error, []string) {

	return nil, c.Destinations[0].To
}

func handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var (
		err           error
		dist          []string
		params        url.Values
		config        DistConfig
		configContent []byte
	)

	configContent, err = readFile("/config.json")

	if err != nil {
		return respond(Response{
			Status: 403,
			Body: JsonBody{
				Type:    "config",
				Success: false,
				Message: "Invalid config path",
			},
		})
	}

	if json.Unmarshal(configContent, &config) != nil {
		return respond(Response{
			Status: 403,
			Body: JsonBody{
				Type:    "config",
				Success: false,
				Message: "Invalid json config",
			},
		})
	}

	configContent = nil

	params, err = url.ParseQuery(request.Body)

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

	if err, dist = resolveDist(config); err != nil {
		return respond(Response{
			Status: 403,
			Body: JsonBody{
				Type:    "host",
				Success: false,
				Message: err.Error(),
			},
		})
	}

	err = utils.Send(
		utils.Email{
			To:   dist,
			Body: "message from: " + params.Get("email"),
		},
		utils.SmtpConfig{
			Host:     os.Getenv("EMAIL_HOST"),
			Port:     os.Getenv("EMAIL_PORT"),
			From:     os.Getenv("EMAIL_FROM"),
			Username: os.Getenv("EMAIL_USERNAME"),
			Password: os.Getenv("EMAIL_PASSWORD"),
		},
	)

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
