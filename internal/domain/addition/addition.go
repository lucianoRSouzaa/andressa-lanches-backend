package addition

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrAdditionNameRequired  = errors.New("o nome do acréscimo é obrigatório")
	ErrAdditionPriceRequired = errors.New("o preço do acréscimo deve ser não negativo")
	ErrAdditionIdInvalid     = errors.New("ID do acréscimo inválido")
	ErrAdditionNotFound      = errors.New("acréscimo não encontrado")
	ErrAdditionIdMandatory   = errors.New("ID do acréscimo é obrigatório")
)

type Addition struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}

func (a *Addition) Validate() error {

	if a.Name == "" {
		return ErrAdditionNameRequired
	}
	if a.Price < 0 {
		return ErrAdditionPriceRequired
	}
	return nil
}
