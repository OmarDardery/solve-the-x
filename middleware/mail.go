package middleware

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(to, code string) error {
	host := os.Getenv("EMAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	user := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")

	msg := gomail.NewMessage()
	msg.SetHeader("From", user)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Verify your account")
	msg.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

	dialer := gomail.NewDialer(host, port, user, pass)
	return dialer.DialAndSend(msg)
}

type Reciever interface {
	GetEmail() string
}

func SendNotification(to, subject, content string) error {
	host := os.Getenv("EMAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	user := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")

	msg := gomail.NewMessage()
	msg.SetHeader("From", user)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", content)

	dialer := gomail.NewDialer(host, port, user, pass)
	return dialer.DialAndSend(msg)
}
