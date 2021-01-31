package requests

// AddToCartRequest request data type for adding item to cart
type AddToCartRequest struct {
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
