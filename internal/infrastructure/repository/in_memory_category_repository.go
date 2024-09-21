package repository

import (
	"context"
	"errors"
	"sync"

	"andressa-lanches/internal/domain/category"

	"github.com/google/uuid"
)

type InMemoryCategoryRepository struct {
	mu         sync.RWMutex
	categories map[uuid.UUID]*category.Category
}

func NewInMemoryCategoryRepository() *InMemoryCategoryRepository {
	return &InMemoryCategoryRepository{
		categories: make(map[uuid.UUID]*category.Category),
	}
}

func (repo *InMemoryCategoryRepository) Create(ctx context.Context, c *category.Category) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	repo.categories[c.ID] = c
	return nil
}

func (repo *InMemoryCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if c, exists := repo.categories[id]; exists {
		return c, nil
	}
	return nil, errors.New("category not found")
}

func (repo *InMemoryCategoryRepository) Update(ctx context.Context, c *category.Category) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.categories[c.ID]; exists {
		repo.categories[c.ID] = c
		return nil
	}
	return errors.New("category not found")
}

func (repo *InMemoryCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.categories[id]; exists {
		delete(repo.categories, id)
		return nil
	}
	return errors.New("category not found")
}

func (repo *InMemoryCategoryRepository) List(ctx context.Context) ([]*category.Category, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	categories := make([]*category.Category, 0, len(repo.categories))
	for _, c := range repo.categories {
		categories = append(categories, c)
	}
	return categories, nil
}
