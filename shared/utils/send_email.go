package utils

import (
	"context"
	"fmt"

	sib_api_v3_sdk "github.com/sendinblue/APIv3-go-library/v2/lib"
)

func SendEmail(
	apiKey string,
	senderEmail string,
	senderName string,
	recipientEmail string,
	recipientName string,
	subject string,
	textContent string,
	html string,
) (string, error) {
	var ctx context.Context

	cfg := sib_api_v3_sdk.NewConfiguration()

	//Configure API key authorization: api-key
	cfg.AddDefaultHeader("api-key", apiKey)

	sib := sib_api_v3_sdk.NewAPIClient(cfg)
	body := sib_api_v3_sdk.SendSmtpEmail{
		HtmlContent: html,
		Subject:     subject,
		Sender: &sib_api_v3_sdk.SendSmtpEmailSender{
			Name:  senderName,
			Email: senderEmail,
		},
		To: []sib_api_v3_sdk.SendSmtpEmailTo{{Name: recipientName, Email: recipientEmail}},
	}
	obj, resp, err := sib.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		fmt.Println(obj)
		fmt.Println(resp)
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error())
		return "", err
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj)
	return obj.MessageId, nil
}
