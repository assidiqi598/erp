package utils

import (
	"gopkg.in/gomail.v2"
)

func SendEmail(smtpKey string, smtpFrom string, senderEmail string, senderName string, recipientEmail string, recipientName string, subject string, textContent string, html string) (string, error) {

	msg := gomail.NewMessage()
	msg.SetHeader("From", senderName+" <"+senderEmail+">")
	msg.SetHeader("To", recipientName+" <"+recipientEmail+">")
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", html)

	// n := gomail.NewDialer()
	return "", nil
}
