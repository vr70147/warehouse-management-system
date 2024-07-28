package utils

import (
	"fmt"
	"log"
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

type NotificationService struct {
	emailSender EmailSender
}

func NewNotificationService(emailSender EmailSender) *NotificationService {
	return &NotificationService{emailSender: emailSender}
}

func (ns *NotificationService) sendNotification(to, subject, body string) error {
	err := ns.emailSender.SendEmail(to, subject, body)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}

func (ns *NotificationService) SendOrderCancellationNotification(email string, orderID uint) error {
	subject := "Your Order Has Been Cancelled"
	body := fmt.Sprintf("Your order with Order ID %d has been cancelled.", orderID)
	return ns.sendNotification(email, subject, body)
}

func (ns *NotificationService) SendOrderCompletionNotification(email string) error {
	return ns.sendNotification(email, "Your Order is Complete", "Dear User, Your order has been completed successfully.")
}
