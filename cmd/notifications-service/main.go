package main

import (
	"log"
	"warehouse-management-system/internal/eventbus"
	"warehouse-management-system/internal/notifications"
)

func main() {
	eb := eventbus.NewEventBus("localhost:6379")
	ns := notifications.NewNotificationsService()

	eb.Subscribe("order_completed", ns.HandleOrderCompleted)
	eb.Subscribe("low_stock_warning", ns.HandleLowStockWarning)
	eb.Subscribe("order_cancelled", ns.HandleOrderCancelled)
	eb.Subscribe("order_shipped", ns.HandleOrderShipped)

	log.Println("Notifications service started")

	select {}
}
