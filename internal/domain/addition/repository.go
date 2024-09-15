package addition

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, addition *Addition) error
	GetByID(ctx context.Context, id uuid.UUID) (*Addition, error)
	Update(ctx context.Context, addition *Addition) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*Addition, error)
}
