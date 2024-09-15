package product

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description,omitempty"`
	CategoryID  uuid.UUID `json:"category_id"`
}
