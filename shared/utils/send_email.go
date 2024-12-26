package utils

import (
	"context"
	"fmt"
	"io"
	"log"

	brevo "github.com/getbrevo/brevo-go/lib"
)

func SendEmail(apiKey string, senderEmail string, senderName string, recipientEmail string, recipientName string, subject string, textContent string, html string) (*brevo.CreateSmtpEmail, error) {

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

	// Log HTTP response status
	if httpResp != nil {
		log.Printf("HTTP Response Status: %v\n", httpResp.Status)
	}

	// Handle errors
	if err != nil {
		if httpResp != nil && httpResp.Body != nil {
			// Safely read the response body
			body, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				log.Printf("Error Response Body: %s\n", string(body))
			} else {
				log.Printf("Error reading response body: %v\n", readErr)
			}
		} else {
			log.Printf("No HTTP response body available.")
		}
		log.Printf("Error sending email: %v\n", err)
		return nil, fmt.Errorf("error sending email: %w", err)
	}

	// Log successful response
	log.Printf("Email sent successfully!\nBrevo Response: %+v\n", resp)

	return &resp, nil
}
