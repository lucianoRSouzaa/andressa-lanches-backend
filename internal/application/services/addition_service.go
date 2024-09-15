package services

import (
	"andressa-lanches/internal/domain/addition"
	"context"
	"errors"

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
	if a.Name == "" {
		return errors.New("o nome do acréscimo é obrigatório")
	}
	if a.Price < 0 {
		return errors.New("o preço do acréscimo deve ser não negativo")
	}

	return s.additionRepo.Create(ctx, a)
}

func (s *additionService) GetAdditionByID(ctx context.Context, id uuid.UUID) (*addition.Addition, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID do acréscimo inválido")
	}

	addition, err := s.additionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if addition == nil {
		return nil, errors.New("acréscimo não encontrado")
	}

	return addition, nil
}

func (s *additionService) UpdateAddition(ctx context.Context, a *addition.Addition) error {
	if a.ID == uuid.Nil {
		return errors.New("ID do acréscimo é obrigatório")
	}
	if a.Name == "" {
		return errors.New("o nome do acréscimo é obrigatório")
	}
	if a.Price < 0 {
		return errors.New("o preço do acréscimo deve ser não negativo")
	}

	existingAddition, err := s.additionRepo.GetByID(ctx, a.ID)
	if err != nil {
		return err
	}
	if existingAddition == nil {
		return errors.New("acréscimo não encontrado")
	}

	return s.additionRepo.Update(ctx, a)
}

func (s *additionService) DeleteAddition(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID do acréscimo inválido")
	}

	existingAddition, err := s.additionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingAddition == nil {
		return errors.New("acréscimo não encontrado")
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
