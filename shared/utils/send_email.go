package utils

import (
	"context"
	"log"

	brevo "github.com/getbrevo/brevo-go/lib"
)

func SendEmail(apiKey string, senderEmail string, senderName string, recipientEmail string, recipientName string, subject string, html string) brevo.CreateSmtpEmail {

	// Create a new Brevo API client
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	client := brevo.NewAPIClient(cfg)

	// Create the transactional email object
	email := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Email: senderEmail, // Replace with your sender email
			Name:  senderName,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: recipientEmail, // Replace with recipient email
				Name:  recipientName,
			},
		},
		Subject:     subject,
		HtmlContent: html,
	}

	// Send the email
	ctx := context.Background()
	resp, httpResp, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, email)

	// Handle errors
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
	} else {
		log.Printf("Email sent successfully!\n")
	}

	log.Printf("HTTP Response Status: %s\n", httpResp.Status)
	log.Printf("Brevo Response: %+v\n", resp)

	return resp
}
