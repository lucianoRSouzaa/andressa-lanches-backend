package services

import (
	"andressa-lanches/internal/domain/category"
	"context"
	"errors"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, c *category.Category) error
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*category.Category, error)
	UpdateCategory(ctx context.Context, c *category.Category) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context) ([]*category.Category, error)
}

type categoryService struct {
	categoryRepo category.Repository
}

func NewCategoryService(categoryRepo category.Repository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, c *category.Category) error {
	if c.Name == "" {
		return errors.New("o nome da categoria é obrigatório")
	}

	return s.categoryRepo.Create(ctx, c)
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID da categoria inválido")
	}

	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("categoria não encontrada")
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, c *category.Category) error {
	if c.ID == uuid.Nil {
		return errors.New("ID da categoria é obrigatório")
	}
	if c.Name == "" {
		return errors.New("o nome da categoria é obrigatório")
	}

	existingCategory, err := s.categoryRepo.GetByID(ctx, c.ID)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return errors.New("categoria não encontrada")
	}

	return s.categoryRepo.Update(ctx, c)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID da categoria inválido")
	}

	existingCategory, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return errors.New("categoria não encontrada")
	}

	return s.categoryRepo.Delete(ctx, id)
}

func (s *categoryService) ListCategories(ctx context.Context) ([]*category.Category, error) {
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
