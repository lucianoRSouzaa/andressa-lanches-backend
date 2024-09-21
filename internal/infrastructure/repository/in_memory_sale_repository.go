package repository

import (
	"context"
	"errors"
	"sync"

	"andressa-lanches/internal/domain/sale"

	"github.com/google/uuid"
)

type InMemorySaleRepository struct {
	mu    sync.RWMutex
	sales map[uuid.UUID]*sale.Sale
}

func NewInMemorySaleRepository() *InMemorySaleRepository {
	return &InMemorySaleRepository{
		sales: make(map[uuid.UUID]*sale.Sale),
	}
}

func (repo *InMemorySaleRepository) Create(ctx context.Context, s *sale.Sale) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	repo.sales[s.ID] = s
	return nil
}

func (repo *InMemorySaleRepository) GetByID(ctx context.Context, id uuid.UUID) (*sale.Sale, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if s, exists := repo.sales[id]; exists {
		return s, nil
	}
	return nil, nil
}

func (repo *InMemorySaleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.sales[id]; exists {
		delete(repo.sales, id)
		return nil
	}
	return errors.New("sale not found")
}

func (repo *InMemorySaleRepository) List(ctx context.Context) ([]*sale.Sale, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	sales := make([]*sale.Sale, 0, len(repo.sales))
	for _, s := range repo.sales {
		sales = append(sales, s)
	}
	return sales, nil
}
