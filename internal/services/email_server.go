package services

import (
	"CardFlow/internal/config"
	"fmt"
	"net/smtp"
)


func SendEmail(data map[string] string) error {
	email := data["Email"]
	firstname := data["FirstName"]
	last_four := data["LastFour"]
	status := data["Status"]
	// Gmail SMTP server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := config.AppEmail
	senderPassword := config.AppPassword

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	switch status {
	case "expired":
		subject := "Your Card Has Expired"
		body := fmt.Sprintf("Dear %s, your card ending with %s has expired.", firstname, last_four)
		message := []byte("Subject: " + subject + "\r\n" +
			"To: " + email + "\r\n" +
			"From: " + senderEmail + "\r\n" +
			"\r\n" +
			body + "\r\n")

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, message)
		return err
	case "expiring":
		subject := "Your Card will soon Expire"
		body := fmt.Sprintf("Dear %s, your card ending with %s is expiring soon.", firstname, last_four)
		message := []byte("Subject: " + subject + "\r\n" +
			"To: " + email + "\r\n" +
			"From: " + senderEmail + "\r\n" +
			"\r\n" +
			body + "\r\n")

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, message)
		return err
	}
	return nil
}

func SendCardTopUpEmail(data map[string]string) error{
	email := data["Email"]
	amount := data["amount"]
	fee := data["fee"]
	firstname := data["FirstName"]
	last_four := data["LastFour"]
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := config.AppEmail
	senderPassword := config.AppPassword
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	subject := "Your Card Was Funded"
		body := fmt.Sprintf("Dear %s, your card ending with %s has been funded with %s.\n  Fee Charged: %s", firstname, last_four, amount, fee)
		message := []byte("Subject: " + subject + "\r\n" +
			"To: " + email + "\r\n" +
			"From: " + senderEmail + "\r\n" +
			"\r\n" +
			body + "\r\n")

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, message)
		return err
}