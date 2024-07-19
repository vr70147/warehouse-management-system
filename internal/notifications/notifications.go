package notifications

import (
	"fmt"
	"log"
	"warehouse-management-system/internal/events"
)

type NotificationService struct{}

func NewNotificationsService() *NotificationService {
	return &NotificationService{}
}

func (ns *NotificationService) HandleOrderCompleted(event *events.Event) {
	userID := event.Data["user_id"].(string)
	orderID := event.Data["order_id"].(string)
	message := fmt.Sprintf("Order %s has been completed", orderID)
	ns.sendNotification(userID, "order Completed", message)
}

func (ns *NotificationService) HandleLowStockWarning(event *events.Event) {
	productID := event.Data["product_id"].(string)
	message := fmt.Sprintf("Product %s is running low in stock", productID)
	ns.sendNotification("admin", "Low Stock Warning", message)
}

func (ns *NotificationService) HandleOrderCancelled(event *events.Event) {
	userID := event.Data["user_id"].(string)
	orderID := event.Data["order_id"].(string)
	message := fmt.Sprintf("Order %s has been cancelled", orderID)
	ns.sendNotification(userID, "order Cancelled", message)
}

func (ns *NotificationService) HandleOrderShipped(event *events.Event) {
	userID := event.Data["user_id"].(string)
	orderID := event.Data["order_id"].(string)
	message := fmt.Sprintf("Order %s has been shipped", orderID)
	ns.sendNotification(userID, "order Shipped", message)
}

func (ns *NotificationService) sendNotification(userID, subject, message string) {
	log.Printf("Sending notification to %s: [%s] %s", userID, subject, message)

	// TODO: send notification
}
