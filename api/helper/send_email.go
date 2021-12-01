package helper

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmailVerCode(emailAddress string, code int) error {
	auth := smtp.PlainAuth("", os.Getenv("MAILER_USERNAME"), os.Getenv("MAILER_PASSWORD"), os.Getenv("MAILER_HOST"))

	from := "userland"
	to := []string{emailAddress}
	msg := []byte(fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", emailAddress) +
		"Subject: Email verification code\r\n" +
		"\r\n" +
		fmt.Sprintf("Here's you code %d it's valid for 60 seconds\r\n", code))

	err := smtp.SendMail(fmt.Sprintf("%s:%s", os.Getenv("MAILER_HOST"), os.Getenv("MAILER_PORT")), auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func SendEmailResetPwdCode(emailAddress string, code int) error {
	auth := smtp.PlainAuth("", os.Getenv("MAILER_USERNAME"), os.Getenv("MAILER_PASSWORD"), os.Getenv("MAILER_HOST"))

	from := "userland"
	to := []string{emailAddress}
	msg := []byte(
		"From : moris@gmail.com\r\n" +
			fmt.Sprintf("To: %s\r\n", emailAddress) +
			"Subject: Email reset password code\r\n" +
			"\r\n" +
			fmt.Sprintf("Here's you code %d it's valid for 60 seconds\r\n", code))

	err := smtp.SendMail(fmt.Sprintf("%s:%s", os.Getenv("MAILER_HOST"), os.Getenv("MAILER_PORT")), auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}
