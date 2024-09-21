package repository

import (
	"context"
	"errors"
	"sync"

	"andressa-lanches/internal/domain/product"

	"github.com/google/uuid"
)

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[uuid.UUID]*product.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[uuid.UUID]*product.Product),
	}
}

func (repo *InMemoryProductRepository) Create(ctx context.Context, p *product.Product) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	repo.products[p.ID] = p
	return nil
}

func (repo *InMemoryProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if p, exists := repo.products[id]; exists {
		return p, nil
	}
	return nil, nil
}

func (repo *InMemoryProductRepository) Update(ctx context.Context, p *product.Product) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.products[p.ID]; exists {
		repo.products[p.ID] = p
		return nil
	}
	return errors.New("product not found")
}

func (repo *InMemoryProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.products[id]; exists {
		delete(repo.products, id)
		return nil
	}
	return errors.New("product not found")
}

func (repo *InMemoryProductRepository) List(ctx context.Context) ([]*product.Product, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	products := make([]*product.Product, 0, len(repo.products))
	for _, p := range repo.products {
		products = append(products, p)
	}
	return products, nil
}
