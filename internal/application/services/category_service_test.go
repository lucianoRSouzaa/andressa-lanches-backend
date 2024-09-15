package services

import (
	"andressa-lanches/internal/domain/category"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, c *category.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	args := m.Called(ctx, id)
	c := args.Get(0)
	if c == nil {
		return nil, args.Error(1)
	}
	return c.(*category.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, c *category.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) List(ctx context.Context) ([]*category.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*category.Category), args.Error(1)
}

func TestCategoryService_CreateCategory_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	testCategory := &category.Category{
		Name:        "Bebidas",
		Description: "Bebidas geladas",
	}

	mockRepo.On("Create", ctx, testCategory).Return(nil)

	err := service.CreateCategory(ctx, testCategory)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_CreateCategory_InvalidName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	testCategory := &category.Category{
		Name: "",
	}

	err := service.CreateCategory(ctx, testCategory)

	assert.Error(t, err)
	assert.Equal(t, "o nome da categoria é obrigatório", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCategoryService_GetCategoryByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	expectedCategory := &category.Category{
		ID:   categoryID,
		Name: "Bebidas",
	}

	mockRepo.On("GetByID", ctx, categoryID).Return(expectedCategory, nil)

	result, err := service.GetCategoryByID(ctx, categoryID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, result)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_UpdateCategory_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	updatedCategory := &category.Category{
		ID:   categoryID,
		Name: "Bebidas Alcoólicas",
	}

	mockRepo.On("GetByID", ctx, categoryID).Return(updatedCategory, nil)
	mockRepo.On("Update", ctx, updatedCategory).Return(nil)

	err := service.UpdateCategory(ctx, updatedCategory)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_DeleteCategory_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()

	mockRepo.On("GetByID", ctx, categoryID).Return(&category.Category{}, nil)
	mockRepo.On("Delete", ctx, categoryID).Return(nil)

	err := service.DeleteCategory(ctx, categoryID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_ListCategories_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expectedCategories := []*category.Category{
		{
			ID:   uuid.New(),
			Name: "Bebidas",
		},
		{
			ID:   uuid.New(),
			Name: "Lanches",
		},
	}

	mockRepo.On("List", ctx).Return(expectedCategories, nil)

	result, err := service.ListCategories(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, result)
	mockRepo.AssertExpectations(t)
}
