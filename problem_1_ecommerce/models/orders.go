package models

const (
	OrderStatusUnprocessed = 0
	OrderStatusCompleted   = 1
	OrderStatusCancel      = 2
)

// Orders data model for orders
type Orders struct {
	Model
	UserID uint        `gorm:"not null" json:"user_id"`
	Status int         `json:"status"` // status: 0 -> unprocessed, 1 -> completed, 2 -> cancel
	Items  []LineItems `gorm:"foreignkey:OrderID" json:"items"`
}
