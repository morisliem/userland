package helper

import (
	"fmt"
	"net/smtp"
)

func SendEmailVerCode(emailAddress string, code int) error {
	auth := smtp.PlainAuth("", "91af5319a267d0", "73d83ac4eef862", "smtp.mailtrap.io")

	to := []string{emailAddress}
	msg := []byte(
		"From : moris@gmail.com\r\n" +
			fmt.Sprintf("To: %s\r\n", emailAddress) +
			"Subject: Email verification code\r\n" +
			"\r\n" +
			fmt.Sprintf("Here's you code %d it's valid for 60 seconds\r\n", code))

	err := smtp.SendMail("smtp.mailtrap.io:2525", auth, "moris@gmail.com", to, msg)
	if err != nil {
		return err
	}
	return nil
}

func SendEmailResetPwdCode(emailAddress string, code int) error {
	auth := smtp.PlainAuth("", "91af5319a267d0", "73d83ac4eef862", "smtp.mailtrap.io")

	to := []string{emailAddress}
	msg := []byte(
		"From : moris@gmail.com\r\n" +
			fmt.Sprintf("To: %s\r\n", emailAddress) +
			"Subject: Email reset password code\r\n" +
			"\r\n" +
			fmt.Sprintf("Here's you code %d it's valid for 60 seconds\r\n", code))

	err := smtp.SendMail("smtp.mailtrap.io:2525", auth, "moris@gmail.com", to, msg)
	if err != nil {
		return err
	}
	return nil
}
