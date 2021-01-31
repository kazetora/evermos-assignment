package models

// Inventories data model for inventory
type Inventories struct {
	Model
	ProductID int `gorm:"not null" json:"product_id"`
	Quantity  int `json:"quantity"`
}
