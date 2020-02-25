package utils

import (
	"fmt"
	"net/smtp"
)

type RawEmail struct {
	From string
	Name string
	Subject string
	Phone string
	Message string
}

type Email struct {
	To string
	Body string
}

type SmtpConfig struct {
	Host string "`json:host`"
	Port string "`json:port`"
	From string "`json:from.email`"
	Name string "`json:from.name`"
	Username  string "`json.username`"
	Password string "`json.password`"
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
