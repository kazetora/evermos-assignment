package models

// Products data model for products
type Products struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
}
