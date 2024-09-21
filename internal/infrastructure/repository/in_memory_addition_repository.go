package repository

import (
	"context"
	"errors"
	"sync"

	"andressa-lanches/internal/domain/addition"

	"github.com/google/uuid"
)

type InMemoryAdditionRepository struct {
	mu        sync.RWMutex
	additions map[uuid.UUID]*addition.Addition
}

func NewInMemoryAdditionRepository() *InMemoryAdditionRepository {
	return &InMemoryAdditionRepository{
		additions: make(map[uuid.UUID]*addition.Addition),
	}
}

func (repo *InMemoryAdditionRepository) Create(ctx context.Context, a *addition.Addition) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	repo.additions[a.ID] = a
	return nil
}

func (repo *InMemoryAdditionRepository) GetByID(ctx context.Context, id uuid.UUID) (*addition.Addition, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if a, exists := repo.additions[id]; exists {
		return a, nil
	}
	return nil, errors.New("addition not found")
}

func (repo *InMemoryAdditionRepository) Update(ctx context.Context, a *addition.Addition) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.additions[a.ID]; exists {
		repo.additions[a.ID] = a
		return nil
	}
	return errors.New("addition not found")
}

func (repo *InMemoryAdditionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.additions[id]; exists {
		delete(repo.additions, id)
		return nil
	}
	return errors.New("addition not found")
}

func (repo *InMemoryAdditionRepository) List(ctx context.Context) ([]*addition.Addition, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	additions := make([]*addition.Addition, 0, len(repo.additions))
	for _, a := range repo.additions {
		additions = append(additions, a)
	}
	return additions, nil
}
