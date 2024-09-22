package services

import (
	"andressa-lanches/internal/domain/addition"
	"context"

	"github.com/google/uuid"
)

type AdditionService interface {
	CreateAddition(ctx context.Context, a *addition.Addition) error
	GetAdditionByID(ctx context.Context, id uuid.UUID) (*addition.Addition, error)
	UpdateAddition(ctx context.Context, a *addition.Addition) error
	DeleteAddition(ctx context.Context, id uuid.UUID) error
	ListAdditions(ctx context.Context) ([]*addition.Addition, error)
}

type additionService struct {
	additionRepo addition.Repository
}

func NewAdditionService(additionRepo addition.Repository) AdditionService {
	return &additionService{
		additionRepo: additionRepo,
	}
}

func (s *additionService) CreateAddition(ctx context.Context, a *addition.Addition) error {
	if err := a.Validate(); err != nil {
		return err
	}

	return s.additionRepo.Create(ctx, a)
}

func (s *additionService) GetAdditionByID(ctx context.Context, id uuid.UUID) (*addition.Addition, error) {
	if id == uuid.Nil {
		return nil, addition.ErrAdditionIdInvalid
	}

	a, err := s.additionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, addition.ErrAdditionNotFound
	}

	return a, nil
}

func (s *additionService) UpdateAddition(ctx context.Context, a *addition.Addition) error {
	if a.ID == uuid.Nil {
		return addition.ErrAdditionIdMandatory
	}

	if err := a.Validate(); err != nil {
		return err
	}

	existingAddition, err := s.additionRepo.GetByID(ctx, a.ID)
	if err != nil {
		return err
	}
	if existingAddition == nil {
		return addition.ErrAdditionNotFound
	}

	return s.additionRepo.Update(ctx, a)
}

func (s *additionService) DeleteAddition(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return addition.ErrAdditionIdInvalid
	}

	existingAddition, err := s.additionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingAddition == nil {
		return addition.ErrAdditionNotFound
	}

	return s.additionRepo.Delete(ctx, id)
}

func (s *additionService) ListAdditions(ctx context.Context) ([]*addition.Addition, error) {
	additions, err := s.additionRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return additions, nil
}
