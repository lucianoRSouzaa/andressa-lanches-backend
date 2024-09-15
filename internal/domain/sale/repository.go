package sale

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, sale *Sale) error
	GetByID(ctx context.Context, id uuid.UUID) (*Sale, error)
	List(ctx context.Context) ([]*Sale, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
