package services

import (
	"andressa-lanches/internal/domain/product"
	"context"
	"errors"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, p *product.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*product.Product, error)
	UpdateProduct(ctx context.Context, p *product.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListProducts(ctx context.Context) ([]*product.Product, error)
}

type productService struct {
	productRepo product.Repository
}

func NewProductService(productRepo product.Repository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, p *product.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return s.productRepo.Create(ctx, p)
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID do produto inválido")
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("produto não encontrado")
	}

	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, p *product.Product) error {
	if p.ID == uuid.Nil {
		return errors.New("ID do produto é obrigatório")
	}
	if p.Name == "" {
		return errors.New("o nome do produto é obrigatório")
	}
	if p.Price <= 0 {
		return errors.New("o preço do produto deve ser positivo")
	}

	existingProduct, err := s.productRepo.GetByID(ctx, p.ID)
	if err != nil {
		return err
	}
	if existingProduct == nil {
		return errors.New("produto não encontrado")
	}

	return s.productRepo.Update(ctx, p)
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID do produto inválido")
	}

	existingProduct, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingProduct == nil {
		return errors.New("produto não encontrado")
	}

	return s.productRepo.Delete(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context) ([]*product.Product, error) {
	products, err := s.productRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
