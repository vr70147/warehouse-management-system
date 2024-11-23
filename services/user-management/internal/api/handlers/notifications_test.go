package handlers

import (
	"testing"
	"user-management/internal/model"
	"user-management/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the SendEmail function
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

// Test for User Registration Notification
func TestSendUserRegistrationNotification(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	userEmail := "test@user.com"
	emailSubject := "Welcome to Our Service"
	emailBody := "Dear User, Welcome to our service!"

	mockEmailSender.On("SendEmail", userEmail, emailSubject, emailBody).Return(nil)

	err := ns.SendUserRegistrationNotification(userEmail)
	assert.NoError(t, err)
	mockEmailSender.AssertExpectations(t)
}

// Test for Order Completion Notification
func TestSendOrderCompletionNotification(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	userEmail := "test@user.com"
	emailSubject := "Your Order is Complete"
	emailBody := "Dear User, Your order has been completed successfully."

	mockEmailSender.On("SendEmail", userEmail, emailSubject, emailBody).Return(nil)

	err := ns.SendOrderCompletionNotification(userEmail)
	assert.NoError(t, err)
	mockEmailSender.AssertExpectations(t)
}

// Test for Low Stock Notification
func TestSendLowStockNotification(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	stakeholderEmail := "stock@warehouse.com"
	emailSubject := "Low Stock Alert"
	emailBody := "Attention: The stock for some items is running low."

	mockEmailSender.On("SendEmail", stakeholderEmail, emailSubject, emailBody).Return(nil)

	err := ns.SendLowStockNotification(stakeholderEmail)
	assert.NoError(t, err)
	mockEmailSender.AssertExpectations(t)
}

// Test for Order Cancellation Notification
func TestSendOrderCancellationNotification(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	userEmail := "test@user.com"
	emailSubject := "Your Order is Canceled"
	emailBody := "Dear User, Your order has been canceled."

	mockEmailSender.On("SendEmail", userEmail, emailSubject, emailBody).Return(nil)

	err := ns.SendOrderCancellationNotification(userEmail)
	assert.NoError(t, err)
	mockEmailSender.AssertExpectations(t)
}

// Test for Password Change Notification
func TestSendPasswordChangeNotification(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	userEmail := "test@user.com"
	emailSubject := "Password Changed Successfully"
	emailBody := "Dear User, Your password has been changed successfully."

	mockEmailSender.On("SendEmail", userEmail, emailSubject, emailBody).Return(nil)

	err := ns.SendPasswordChangeNotification(userEmail)
	assert.NoError(t, err)
	mockEmailSender.AssertExpectations(t)
}

// Test for Failed Login Attempt Notification
func TestSendFailedLoginAttemptNotification(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	userEmail := "test@user.com"
	emailSubject := "Failed Login Attempt"
	emailBody := "Dear User, There was a failed login attempt on your account."

	mockEmailSender.On("SendEmail", userEmail, emailSubject, emailBody).Return(nil)

	err := ns.SendFailedLoginAttemptNotification(userEmail)
	assert.NoError(t, err)
	mockEmailSender.AssertExpectations(t)
}

// Test for Monthly Account Activity Summary Notification
func TestSendMonthlySummaryNotification(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	user := model.User{
		Email: "test@user.com",
		Name:  "User",
	}

	emailSubject := "Monthly Account Activity Summary"
	summary := "Here is your account activity summary for the last month."
	emailBody := "Dear User,\n\nHere is your account activity summary for the last month:\n\n" + summary

	mockEmailSender.On("SendEmail", user.Email, emailSubject, emailBody).Return(nil)

	err := ns.SendMonthlySummaryNotification(user, summary)
	assert.NoError(t, err)
	mockEmailSender.AssertExpectations(t)
}
