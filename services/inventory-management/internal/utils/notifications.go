package utils

import (
	"fmt"
	"log"
)

// NotificationService handles sending notifications
type NotificationService struct {
	emailSender EmailSender
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(emailSender EmailSender) *NotificationService {
	return &NotificationService{emailSender: emailSender}
}

func (ns *NotificationService) sendNotification(to, subject, body string) error {
	if ns.emailSender == nil {
		return fmt.Errorf("email sender is not initialized")
	}
	log.Printf("Sending notification to %s: [%s] %s", to, subject, body)
	return ns.emailSender.SendEmail(to, subject, body)
}

func (ns *NotificationService) SendLowStockNotification(email string) error {
	return ns.sendNotification(email, "Low Stock Alert", "Attention: The stock for some items is running low.")
}
