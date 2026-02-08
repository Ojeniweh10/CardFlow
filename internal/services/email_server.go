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
	email := data["email"]
	amount := data["amount"]
	fee := data["fee"]
	firstname := data["firstName"]
	last_four := data["lastFour"]
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

func SendCardDebitEmail(data map[string]string) error{
	email := data["email"]
	amount := data["amount"]
	fee := data["fee"]
	firstname := data["firstname"]
	last_four := data["lastfour"]
	balance := data["balance"]
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := config.AppEmail
	senderPassword := config.AppPassword
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	subject := "Your Card Was Debited"
		body := fmt.Sprintf("Dear %s, your card ending with %s has been Debited the amount of  %s.\n  Fee Charged: %s.\n Balance: %s ", firstname, last_four, amount, fee, balance)
		message := []byte("Subject: " + subject + "\r\n" +
			"To: " + email + "\r\n" +
			"From: " + senderEmail + "\r\n" +
			"\r\n" +
			body + "\r\n")

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, message)
		return err
}


func SendDebitReversalEmail(data map[string]string) error{
	email := data["email"]
	amount := data["amount"]
	firstname := data["firstname"]
	last_four := data["lastfour"]
	balance := data["balance"]
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := config.AppEmail
	senderPassword := config.AppPassword
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	subject := "Your Transaction was Reversed"
		body := fmt.Sprintf("Dear %s, your card ending with %s which was Debited the amount of  %s.\n Has Been Reversed.\n Balance: %s ", firstname, last_four, amount, balance)
		message := []byte("Subject: " + subject + "\r\n" +
			"To: " + email + "\r\n" +
			"From: " + senderEmail + "\r\n" +
			"\r\n" +
			body + "\r\n")

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, message)
		return err
}

func SendRefundEmail(data map[string]string) error{
	email := data["email"]
	amount := data["amount"]
	firstname := data["firstname"]
	last_four := data["lastfour"]
	balance := data["balance"]
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := config.AppEmail
	senderPassword := config.AppPassword
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	subject := "Your Transaction was Refunded"
		body := fmt.Sprintf("Dear %s, your card ending with %s Has been refunded the amount of  %s.\n Balance: %s ", firstname, last_four, amount, balance)
		message := []byte("Subject: " + subject + "\r\n" +
			"To: " + email + "\r\n" +
			"From: " + senderEmail + "\r\n" +
			"\r\n" +
			body + "\r\n")

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, message)
		return err
}