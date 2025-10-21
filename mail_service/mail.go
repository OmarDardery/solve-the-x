package mail_service

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendVerificationEmail(to, code string) error {
	fromEmail := os.Getenv("SENDER_EMAIL")
	apiKey := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail("Solve-The-X", fromEmail)
	subject := "Verify your account"
	toEmail := mail.NewEmail("", to)
	plainTextContent := fmt.Sprintf("Your verification code is: %s", code)
	htmlContent := fmt.Sprintf("<p>Your verification code is: <b>%s</b></p>", code)

	message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: %d - %s", response.StatusCode, response.Body)
	}

	return nil
}

func SendNotification(to, subject, content string) error {
	fromEmail := os.Getenv("SENDER_EMAIL")
	apiKey := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail("Solve-The-X", fromEmail)
	toEmail := mail.NewEmail("", to)

	message := mail.NewSingleEmail(from, subject, toEmail, content, "<p>"+content+"</p>")
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: %d - %s", response.StatusCode, response.Body)
	}

	return nil
}
