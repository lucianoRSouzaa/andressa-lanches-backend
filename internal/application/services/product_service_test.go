package services

import (
	"andressa-lanches/internal/domain/product"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, p *product.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	args := m.Called(ctx, id)
	p := args.Get(0)
	if p == nil {
		return nil, args.Error(1)
	}
	return p.(*product.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, p *product.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) List(ctx context.Context) ([]*product.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*product.Product), args.Error(1)
}

func TestProductService_CreateProduct_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	testProduct := &product.Product{
		Name:        "Sanduíche",
		Price:       12.50,
		Description: "Delicioso sanduíche",
		CategoryID:  uuid.New(),
	}

	mockRepo.On("Create", ctx, testProduct).Return(nil)

	err := service.CreateProduct(ctx, testProduct)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_InvalidName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	testProduct := &product.Product{
		Name:  "",
		Price: 12.50,
	}

	err := service.CreateProduct(ctx, testProduct)

	assert.Error(t, err)
	assert.Equal(t, "o nome do produto é obrigatório", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestProductService_CreateProduct_InvalidPrice(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	testProduct := &product.Product{
		Name:  "Sanduíche",
		Price: -5.00,
	}

	err := service.CreateProduct(ctx, testProduct)

	assert.Error(t, err)
	assert.Equal(t, "o preço do produto deve ser positivo", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestProductService_GetProductByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	productID := uuid.New()
	expectedProduct := &product.Product{
		ID:    productID,
		Name:  "Sanduíche",
		Price: 12.50,
	}

	mockRepo.On("GetByID", ctx, productID).Return(expectedProduct, nil)

	result, err := service.GetProductByID(ctx, productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, result)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	productID := uuid.New()

	mockRepo.On("GetByID", ctx, productID).Return(nil, nil)

	result, err := service.GetProductByID(ctx, productID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "produto não encontrado", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	productID := uuid.New()
	updatedProduct := &product.Product{
		ID:         productID,
		Name:       "Sanduíche Atualizado",
		Price:      15.00,
		CategoryID: uuid.New(),
	}

	mockRepo.On("GetByID", ctx, productID).Return(updatedProduct, nil)
	mockRepo.On("Update", ctx, updatedProduct).Return(nil)

	err := service.UpdateProduct(ctx, updatedProduct)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	productID := uuid.New()

	mockRepo.On("GetByID", ctx, productID).Return(&product.Product{}, nil)
	mockRepo.On("Delete", ctx, productID).Return(nil)

	err := service.DeleteProduct(ctx, productID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_ListProducts_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProducts := []*product.Product{
		{
			ID:    uuid.New(),
			Name:  "Sanduíche",
			Price: 12.50,
		},
		{
			ID:    uuid.New(),
			Name:  "Suco",
			Price: 5.00,
		},
	}

	mockRepo.On("List", ctx).Return(expectedProducts, nil)

	result, err := service.ListProducts(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, result)
	mockRepo.AssertExpectations(t)
}
