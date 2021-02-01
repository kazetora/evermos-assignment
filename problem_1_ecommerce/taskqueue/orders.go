package taskqueue

import (
	"log"
	"time"

	"github.com/adjust/rmq/v3"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/database"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/models"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/repositories"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/storage"
)

// OrderTask data type for order queue's task
type OrderTask struct {
	UserID        uint
	TransactionID string
	ProductID     uint
	Quantity      int
}

// PublshToOrderTaskQueue publish ask to order task queue
func PublshToOrderTaskQueue(transactionKey string) error {

	err := orderTaskQueue.Publish(transactionKey)
	if err != nil {
		log.Printf("Error adding to task queue: %s", err.Error())
		return err
	}

	return nil

}

func startConsumingOrderTask() error {
	err := orderTaskQueue.StartConsuming(1000, time.Second)
	if err != nil {
		log.Printf("Error start consuming queue: %s", err.Error())
		return err
	}

	_, err = orderTaskQueue.AddConsumerFunc("orderConsumer", func(delivery rmq.Delivery) {
		// var task OrderTask
		transactionKey := delivery.Payload()
		task, err := storage.GetTransaction(transactionKey)
		if err != nil {
			log.Printf("Error retrieving task: %s", err.Error())

			if err = delivery.Reject(); err != nil {
				log.Printf("Error rejecting order task: %s", err.Error())
			}

		} else {
			// check if users currently has an active order
			var activeOrder models.Orders
			db := database.GetDatabase()
			findRes := db.Where("user_id = ? AND status = ?", task.UserID, models.OrderStatusUnprocessed).First(&activeOrder)
			if findRes.RecordNotFound() {
				// Create new order for this user.
				newOrder, err := repositories.CreateNewOrder(db, task.UserID, task.ProductID, task.Quantity)

				if err != nil {
					// update transaction status
					updateTransaction(transactionKey, task, models.TransactionStatusError)
					if err = delivery.Reject(); err != nil {
						log.Printf("Error rejecting order task: %s", err.Error())
					}
					return

				}

				activeOrder = *newOrder

			} else {
				// update line items in order
				existingOrder, err := repositories.UpdateOrderLineItems(db, activeOrder.ID, task.ProductID, task.Quantity)
				if err != nil {

					// update transaction status
					updateTransaction(transactionKey, task, models.TransactionStatusError)

					if err = delivery.Reject(); err != nil {
						log.Printf("Error rejecting order task: %s", err.Error())

					}
					return
				}

				activeOrder = *existingOrder

			}

			// update transaction for success and send Ack for queue
			updateTransaction(transactionKey, task, models.TransactionStatusSuccess)
			if err = delivery.Ack(); err != nil {
				log.Printf("Error acking order task: %s", err.Error())
			}

		}

	})

	return nil
}

func updateTransaction(key string, value models.TransactionCache, status int) {
	value.Status = status
	if err := storage.UpdateTransaction(key, value); err != nil {
		log.Printf("Error updating transaction: %s", err.Error())
	}
}
