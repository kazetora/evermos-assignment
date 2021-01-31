package taskqueue

import (
	"encoding/json"
	"log"
	"time"

	"github.com/adjust/rmq/v3"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/database"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/models"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/repositories"
)

// OrderTask data type for order queue's task
type OrderTask struct {
	UserID        uint
	TransactionID uint
	ProductID     uint
	Quantity      int
}

// PublshToOrderTaskQueue publish ask to order task queue
func PublshToOrderTaskQueue(task OrderTask) error {

	taskBytes, err := json.Marshal(task)
	if err != nil {
		log.Printf("Error adding to task queue: %s", err.Error())
		return err
	}

	err = orderTaskQueue.PublishBytes(taskBytes)
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
		var task OrderTask
		if err := json.Unmarshal([]byte(delivery.Payload()), &task); err != nil {
			log.Printf("Error unmarshal parameter: %s", err.Error())

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
					updateTransaction(task, models.TransactionStatusError)
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

					updateTransaction(task, models.TransactionStatusError)
					if err = delivery.Reject(); err != nil {
						log.Printf("Error rejecting order task: %s", err.Error())

					}
					return
				}

				activeOrder = *existingOrder

			}

			// update transaction for success and send Ack for queue
			updateTransaction(task, models.TransactionStatusSuccess)
			if err = delivery.Ack(); err != nil {
				log.Printf("Error acking order task: %s", err.Error())
			}

		}

	})

	return nil
}

func updateTransaction(task OrderTask, value int) error {
	db := database.GetDatabase()

	if err := db.Model(models.Transactions{}).Where("id = ?", task.TransactionID).Update("status", value).Error; err != nil {
		log.Printf("Error updating transaction: %s", err.Error())
		return err
	}

	return nil

}
