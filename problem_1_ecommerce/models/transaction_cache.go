package models

type TransactionCache struct {
	UserID    uint `json:"user_id"`
	Status    int  `json:"status"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
