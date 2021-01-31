package models

const (
	TransactionStatusUnprocessed = 0
	TransactionStatusSuccess     = 1
	TransactionStatusError       = 2
)

// Transactions data model for transactions
type Transactions struct {
	Model
	UserID  uint `gorm:"not null" json:"user_id"`
	Status  int  `json:"status"` // status: 0 -> unprocessed, 1 -> success, 2 -> cancel
	OrderID uint `json:"order_id"`
}
