package utils

import (
	"os"

	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

// DefaultEmailSender is the default implementation of EmailSender
type DefaultEmailSender struct{}

func (s *DefaultEmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_ADDRESS"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), 587, os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"))
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
