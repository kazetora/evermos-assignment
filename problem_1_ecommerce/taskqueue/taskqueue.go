package taskqueue

import (
	"log"
	"os"

	"github.com/adjust/rmq/v3"
)

const (
	// OrderProcessingQueue task queue name for order processing
	OrderProcessingQueue = "order_task_queue"
)

var connection rmq.Connection
var orderTaskQueue rmq.Queue

// Setup init taskqueues
func Setup() error {
	var err error
	connection, err = rmq.OpenConnection("taskqueue", "tcp", os.Getenv("REDIS_SERVER"), 2, nil)
	if err != nil {
		log.Printf("Error setting task queue: %s", err.Error())
		return err
	}

	orderTaskQueue, err = connection.OpenQueue(OrderProcessingQueue)
	if err != nil {
		log.Printf("Error setting task queue: %s", err.Error())
		return err
	}

	return nil
}

// StartConsumers start consumers in all queues
func StartConsumers() error {
	if err := startConsumingOrderTask(); err != nil {
		return err
	}

	return nil
}

// // GetConnection return the rmq connection
// func GetConnection() rmq.Connection {
// 	return connection
// }

// // GetTaskQueue return the rmq queue
// func GetTaskQueue(name string) (rmq.Queue, error) {
// 	switch name {
// 	case OrderProcessingQueue:
// 		return orderTaskQueue, nil
// 	default:
// 		return nil, errors.New("Task queue does not exist")
// 	}
// }
