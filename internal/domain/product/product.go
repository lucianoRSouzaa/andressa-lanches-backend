package product

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrProductNameRequired  = errors.New("o nome do produto é obrigatório")
	ErrProductPricePositive = errors.New("o preço do produto deve ser positivo")
	ErrProductCategoryID    = errors.New("o ID da categoria do produto é obrigatório")
	ErrProductNotFound      = errors.New("produto não encontrado")
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description,omitempty"`
	CategoryID  uuid.UUID `json:"category_id"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrProductNameRequired
	}
	if p.Price <= 0 {
		return ErrProductPricePositive
	}
	if p.CategoryID == uuid.Nil {
		return ErrProductCategoryID
	}
	return nil
}
