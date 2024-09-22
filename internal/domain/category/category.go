package category

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCategoryNameRequired = errors.New("o nome da categoria é obrigatório")
	ErrCategoryIdInvalid    = errors.New("ID da categoria inválido")
	ErrCategoryIdRequired   = errors.New("ID da categoria é obrigatório")
	ErrCategoryNotFound     = errors.New("categoria não encontrada")
)

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

func (p *Category) Validate() error {
	if p.Name == "" {
		return ErrCategoryNameRequired
	}
	return nil
}
