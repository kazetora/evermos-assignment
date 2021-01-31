package repositories

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/models"
)

func isStockAvailable(db *gorm.DB, productID uint, qty int) (bool, error) {
	var inventory models.Inventories
	if err := db.Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		return false, err
	}
	if inventory.Quantity < qty {
		return false, errors.New("Not enough stock")
	}

	return true, nil
}

func updateStock(db *gorm.DB, productID uint, qty int) error {
	var inventory models.Inventories
	if err := db.Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		return err
	}
	inventory.Quantity = inventory.Quantity - qty

	if err := db.Save(&inventory).Error; err != nil {
		log.Printf("Error updating inventory: %s", err.Error())
		return err
	}

	return nil
}

// CreateNewOrder create new unprocessed order for user
func CreateNewOrder(db *gorm.DB, userID uint, productID uint, qty int) (*models.Orders, error) {
	// check if the inventory can handle the quantity
	if stockAvailable, err := isStockAvailable(db, productID, qty); err != nil || !stockAvailable {
		log.Println(err.Error())
		return nil, errors.New("Cannot continue with the transaction")
	}

	order := models.Orders{
		UserID: userID,
		Items: []models.LineItems{
			models.LineItems{
				ProductID: productID,
				Quantity:  qty,
			},
		},
	}

	createRes := db.Create(&order)
	if createRes.Error != nil {
		return nil, createRes.Error
	}

	// update inventory
	if err := updateStock(db, productID, qty); err != nil {
		return nil, err
	}

	return &order, nil
}

// UpdateOrderLineItems update line item in order.
// If the same product id exists in current line items, update the quantity
func UpdateOrderLineItems(db *gorm.DB, orderID uint, productID uint, qty int) (*models.Orders, error) {

	// check if the inventory can handle the quantity
	if stockAvailable, err := isStockAvailable(db, productID, qty); err != nil || !stockAvailable {
		log.Println(err.Error())
		return nil, errors.New("Cannot continue with the transaction")
	}

	// check if the same item exist in order's line items
	var findItem models.LineItems
	queryRes := db.Where("order_id = ?", orderID).First(&findItem)
	if queryRes.RecordNotFound() {
		// create new line item for the order
		findItem.OrderID = orderID
		findItem.ProductID = productID
		findItem.Quantity = qty
		if err := db.Create(&findItem).Error; err != nil {
			return nil, err
		}
	} else {
		// update the quantity
		findItem.Quantity = qty
		if err := db.Save(&findItem).Error; err != nil {
			return nil, err
		}
	}

	var ret models.Orders
	if err := db.Where("id = ?", orderID).Error; err != nil {
		return nil, err
	}

	// update inventory
	if err := updateStock(db, productID, qty); err != nil {
		return nil, err
	}

	return &ret, nil

}
