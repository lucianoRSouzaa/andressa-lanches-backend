package services

import (
	"andressa-lanches/internal/domain/addition"
	"andressa-lanches/internal/domain/product"
	"andressa-lanches/internal/domain/sale"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type SaleService interface {
	CreateSale(ctx context.Context, s *sale.Sale) error
	GetSaleByID(ctx context.Context, id uuid.UUID) (*sale.Sale, error)
	ListSales(ctx context.Context) ([]*sale.Sale, error)
	DeleteSale(ctx context.Context, id uuid.UUID) error
}

type saleService struct {
	saleRepo     sale.Repository
	productRepo  product.Repository
	additionRepo addition.Repository
}

func NewSaleService(
	saleRepo sale.Repository,
	productRepo product.Repository,
	additionRepo addition.Repository,
) SaleService {
	return &saleService{
		saleRepo:     saleRepo,
		productRepo:  productRepo,
		additionRepo: additionRepo,
	}
}

func (s *saleService) CreateSale(ctx context.Context, sale *sale.Sale) error {
	if sale.Date.IsZero() {
		sale.Date = time.Now()
	}

	var totalSaleAmount float64

	for i := range sale.Items {
		item := &sale.Items[i]

		if item.ProductID == uuid.Nil {
			return errors.New("ID do produto é obrigatório para o item da venda")
		}

		prod, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil || prod == nil {
			return errors.New("produto não encontrado")
		}
		item.UnitPrice = prod.Price

		if item.Quantity <= 0 {
			return errors.New("a quantidade deve ser positiva")
		}

		var totalAdditionsPrice float64
		for j := range item.Additions {
			additionID := item.Additions[j].ID
			if additionID == uuid.Nil {
				return errors.New("ID do acréscimo inválido")
			}

			add, err := s.additionRepo.GetByID(ctx, additionID)
			if err != nil || add == nil {
				return errors.New("acréscimo não encontrado")
			}
			totalAdditionsPrice += add.Price
			item.Additions[j] = *add
		}

		item.TotalPrice = (item.UnitPrice + totalAdditionsPrice) * float64(item.Quantity)
		totalSaleAmount += item.TotalPrice
	}

	sale.TotalAmount = totalSaleAmount - sale.Discount + sale.AdditionalCharges

	return s.saleRepo.Create(ctx, sale)
}

func (s *saleService) GetSaleByID(ctx context.Context, id uuid.UUID) (*sale.Sale, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID da venda inválido")
	}

	sale, err := s.saleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, errors.New("venda não encontrada")
	}

	return sale, nil
}

func (s *saleService) ListSales(ctx context.Context) ([]*sale.Sale, error) {
	sales, err := s.saleRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (s *saleService) DeleteSale(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID da venda inválido")
	}

	existingSale, err := s.saleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingSale == nil {
		return errors.New("venda não encontrada")
	}

	return s.saleRepo.Delete(ctx, id)
}
