package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/adjust/rmq/v3"
	"github.com/jinzhu/gorm"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/helpers"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/models"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/requests"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/taskqueue"
)

// OrderController controller class for Order management
type OrderController struct {
	Db      *gorm.DB
	RmqConn rmq.Connection
}

// NewOrderController OderController class constructor
// params:
// - db: gorm db object
func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{
		Db: db,
	}
}

// AddToCart create new order
func (c *OrderController) AddToCart(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	req := requests.AddToCartRequest{}
	if err := decoder.Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// create a transaction
	transaction := models.Transactions{
		UserID: req.UserID,
		Status: models.TransactionStatusUnprocessed,
	}

	if err := c.Db.Create(&transaction).Error; err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err := taskqueue.PublshToOrderTaskQueue(taskqueue.OrderTask{
		UserID:        req.UserID,
		TransactionID: transaction.ID,
		ProductID:     req.ProductID,
		Quantity:      req.Quantity,
	})

	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, transaction)
}
