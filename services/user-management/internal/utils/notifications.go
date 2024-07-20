package utils

import (
	"fmt"
	"log"
	"user-management/internal/model"
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

func (ns *NotificationService) SendUserRegistrationNotification(email string) error {
	return ns.sendNotification(email, "Welcome to Our Service", "Dear User, Welcome to our service!")
}

func (ns *NotificationService) SendOrderCompletionNotification(email string) error {
	return ns.sendNotification(email, "Your Order is Complete", "Dear User, Your order has been completed successfully.")
}

func (ns *NotificationService) SendLowStockNotification(email string) error {
	return ns.sendNotification(email, "Low Stock Alert", "Attention: The stock for some items is running low.")
}

func (ns *NotificationService) SendOrderCancellationNotification(email string) error {
	return ns.sendNotification(email, "Your Order is Canceled", "Dear User, Your order has been canceled.")
}

func (ns *NotificationService) SendPasswordChangeNotification(email string) error {
	return ns.sendNotification(email, "Password Changed Successfully", "Dear User, Your password has been changed successfully.")
}

func (ns *NotificationService) SendFailedLoginAttemptNotification(email string) error {
	return ns.sendNotification(email, "Failed Login Attempt", "Dear User, There was a failed login attempt on your account.")
}

func (ns *NotificationService) SendMonthlySummaryNotification(user model.User, summary string) error {
	emailSubject := "Monthly Account Activity Summary"
	emailBody := fmt.Sprintf("Dear %s,\n\nHere is your account activity summary for the last month:\n\n%s", user.Name, summary)
	return ns.emailSender.SendEmail(user.Email, emailSubject, emailBody)
}
