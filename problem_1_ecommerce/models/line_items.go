package models

// LineItems data model for order line items
type LineItems struct {
	Model
	OrderID   uint      `gorm:"not null" json:"order_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Product   *Products `gorm:"foreignkey:ProductID"`
}
