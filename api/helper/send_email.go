package helper

import (
	"fmt"
	"net/smtp"
)

func SendEmail(emailAddress string, code int) error {
	auth := smtp.PlainAuth("", "94ca30f7cfa29a", "8fc65ab0d35537", "smtp.mailtrap.io")

	to := []string{emailAddress}
	msg := []byte(fmt.Sprintf("To: %s\r\n", emailAddress) +
		"Subject: Email verification?\r\n" +
		"\r\n" +
		fmt.Sprintf("Here's you code %d \r\n", code))
	err := smtp.SendMail("smtp.mailtrap.io:2525", auth, "moris@gmail.com", to, msg)

	if err != nil {
		return err
	}
	return nil
}
