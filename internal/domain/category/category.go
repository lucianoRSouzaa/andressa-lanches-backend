package category

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCategoryNameRequired = errors.New("o nome da categoria é obrigatório")
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
