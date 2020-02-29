package utils

import (
	"fmt"
	"net/smtp"
)

type RawEmail struct {
	From    string `json:"email"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type Email struct {
	To   string
	Body string
}

type SmtpConfig struct {
	Host     string
	Port     string
	From     string
	Username string
	Password string
}

func Send(email Email, config SmtpConfig) error {

	var mailer string = fmt.Sprintf("%s:%s", config.Host, config.Port)

	return smtp.SendMail(
		mailer,
		smtp.PlainAuth("", config.Username, config.Password, mailer),
		config.From,
		[]string{email.To},
		[]byte(email.Body),
	)
}
