package services

import (
	"andressa-lanches/internal/domain/addition"
	"andressa-lanches/internal/domain/product"
	"andressa-lanches/internal/domain/sale"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks dos repositórios
type MockSaleRepository struct {
	mock.Mock
}

func (m *MockSaleRepository) Create(ctx context.Context, s *sale.Sale) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockSaleRepository) GetByID(ctx context.Context, id uuid.UUID) (*sale.Sale, error) {
	args := m.Called(ctx, id)
	s := args.Get(0)
	if s == nil {
		return nil, args.Error(1)
	}
	return s.(*sale.Sale), args.Error(1)
}

func (m *MockSaleRepository) List(ctx context.Context) ([]*sale.Sale, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*sale.Sale), args.Error(1)
}

func (m *MockSaleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSaleService_CreateSale_Success(t *testing.T) {
	ctx := context.Background()
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	mockAdditionRepo := new(MockAdditionRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo, mockAdditionRepo)

	productID := uuid.New()
	additionID := uuid.New()

	testSale := &sale.Sale{
		Discount:          0.0,
		AdditionalCharges: 0.0,
		Items: []sale.SaleItem{
			{
				ProductID: productID,
				Quantity:  2,
				Additions: []addition.Addition{
					{ID: additionID},
				},
			},
		},
	}

	mockProductRepo.On("GetByID", ctx, productID).Return(&product.Product{
		ID:    productID,
		Name:  "Sanduíche",
		Price: 10.00,
	}, nil)

	mockAdditionRepo.On("GetByID", ctx, additionID).Return(&addition.Addition{
		ID:    additionID,
		Name:  "Bacon",
		Price: 2.50,
	}, nil)

	mockSaleRepo.On("Create", ctx, mock.AnythingOfType("*sale.Sale")).Return(nil)

	err := service.CreateSale(ctx, testSale)

	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
	mockAdditionRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)

	expectedTotalPrice := (10.00 + 2.50) * 2
	assert.Equal(t, expectedTotalPrice, testSale.TotalAmount)
}

func TestSaleService_CreateSale_InvalidProductID(t *testing.T) {
	ctx := context.Background()
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	mockAdditionRepo := new(MockAdditionRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo, mockAdditionRepo)

	testSale := &sale.Sale{
		Items: []sale.SaleItem{
			{
				ProductID: uuid.Nil,
				Quantity:  1,
			},
		},
	}

	err := service.CreateSale(ctx, testSale)

	assert.Error(t, err)
	assert.Equal(t, "ID do produto é obrigatório para o item da venda", err.Error())
	mockProductRepo.AssertNotCalled(t, "GetByID")
	mockSaleRepo.AssertNotCalled(t, "Create")
}

func TestSaleService_GetSaleByID_Success(t *testing.T) {
	ctx := context.Background()
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	mockAdditionRepo := new(MockAdditionRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo, mockAdditionRepo)

	saleID := uuid.New()
	expectedSale := &sale.Sale{
		ID:          saleID,
		Date:        time.Now(),
		TotalAmount: 50.00,
	}

	mockSaleRepo.On("GetByID", ctx, saleID).Return(expectedSale, nil)

	result, err := service.GetSaleByID(ctx, saleID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSale, result)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_DeleteSale_Success(t *testing.T) {
	ctx := context.Background()
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	mockAdditionRepo := new(MockAdditionRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo, mockAdditionRepo)

	saleID := uuid.New()

	mockSaleRepo.On("GetByID", ctx, saleID).Return(&sale.Sale{}, nil)
	mockSaleRepo.On("Delete", ctx, saleID).Return(nil)

	err := service.DeleteSale(ctx, saleID)

	assert.NoError(t, err)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_ListSales_Success(t *testing.T) {
	ctx := context.Background()
	mockSaleRepo := new(MockSaleRepository)
	mockProductRepo := new(MockProductRepository)
	mockAdditionRepo := new(MockAdditionRepository)
	service := NewSaleService(mockSaleRepo, mockProductRepo, mockAdditionRepo)

	expectedSales := []*sale.Sale{
		{
			ID:          uuid.New(),
			Date:        time.Now(),
			TotalAmount: 50.00,
		},
		{
			ID:          uuid.New(),
			Date:        time.Now(),
			TotalAmount: 30.00,
		},
	}

	mockSaleRepo.On("List", ctx).Return(expectedSales, nil)

	result, err := service.ListSales(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedSales, result)
	mockSaleRepo.AssertExpectations(t)
}
