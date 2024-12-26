package utils

import (
	"context"
	"fmt"
	"log"

	brevo "github.com/getbrevo/brevo-go/lib"
)

func SendEmail(apiKey string, senderEmail string, senderName string, recipientEmail string, recipientName string, subject string, textContent string, html string) error {

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
		TextContent: textContent,
		HtmlContent: html,
	}

	// Send the email
	ctx := context.Background()
	resp, httpResp, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, email)

	log.Printf("HTTP Response Status: %v\n", httpResp)
	log.Printf("Brevo Response: %+v\n", resp)

	// Handle errors
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	log.Printf("Email sent successfully!\n")

	return nil
}
