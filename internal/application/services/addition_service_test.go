package services

import (
	"andressa-lanches/internal/domain/addition"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAdditionRepository struct {
	mock.Mock
}

func (m *MockAdditionRepository) Create(ctx context.Context, a *addition.Addition) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockAdditionRepository) GetByID(ctx context.Context, id uuid.UUID) (*addition.Addition, error) {
	args := m.Called(ctx, id)
	add := args.Get(0)
	if add == nil {
		return nil, args.Error(1)
	}
	return add.(*addition.Addition), args.Error(1)
}

func (m *MockAdditionRepository) Update(ctx context.Context, a *addition.Addition) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockAdditionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAdditionRepository) List(ctx context.Context) ([]*addition.Addition, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*addition.Addition), args.Error(1)
}

func TestAdditionService_CreateAddition_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAdditionRepository)
	service := NewAdditionService(mockRepo)

	testAddition := &addition.Addition{
		Name:  "Bacon",
		Price: 2.50,
	}

	mockRepo.On("Create", ctx, testAddition).Return(nil)

	err := service.CreateAddition(ctx, testAddition)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAdditionService_CreateAddition_InvalidName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAdditionRepository)
	service := NewAdditionService(mockRepo)

	testAddition := &addition.Addition{
		Name:  "",
		Price: 2.50,
	}

	err := service.CreateAddition(ctx, testAddition)

	assert.Error(t, err)
	assert.Equal(t, "o nome do acréscimo é obrigatório", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestAdditionService_CreateAddition_InvalidPrice(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAdditionRepository)
	service := NewAdditionService(mockRepo)

	testAddition := &addition.Addition{
		Name:  "Bacon",
		Price: -1.00,
	}

	err := service.CreateAddition(ctx, testAddition)

	assert.Error(t, err)
	assert.Equal(t, "o preço do acréscimo deve ser não negativo", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestAdditionService_GetAdditionByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAdditionRepository)
	service := NewAdditionService(mockRepo)

	additionID := uuid.New()
	expectedAddition := &addition.Addition{
		ID:    additionID,
		Name:  "Bacon",
		Price: 2.50,
	}

	mockRepo.On("GetByID", ctx, additionID).Return(expectedAddition, nil)

	result, err := service.GetAdditionByID(ctx, additionID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAddition, result)
	mockRepo.AssertExpectations(t)
}

func TestAdditionService_UpdateAddition_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAdditionRepository)
	service := NewAdditionService(mockRepo)

	additionID := uuid.New()
	updatedAddition := &addition.Addition{
		ID:    additionID,
		Name:  "Queijo",
		Price: 1.50,
	}

	mockRepo.On("GetByID", ctx, additionID).Return(updatedAddition, nil)
	mockRepo.On("Update", ctx, updatedAddition).Return(nil)

	err := service.UpdateAddition(ctx, updatedAddition)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAdditionService_DeleteAddition_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAdditionRepository)
	service := NewAdditionService(mockRepo)

	additionID := uuid.New()

	mockRepo.On("GetByID", ctx, additionID).Return(&addition.Addition{}, nil)
	mockRepo.On("Delete", ctx, additionID).Return(nil)

	err := service.DeleteAddition(ctx, additionID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAdditionService_ListAdditions_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAdditionRepository)
	service := NewAdditionService(mockRepo)

	expectedAdditions := []*addition.Addition{
		{
			ID:    uuid.New(),
			Name:  "Bacon",
			Price: 2.50,
		},
		{
			ID:    uuid.New(),
			Name:  "Queijo",
			Price: 1.50,
		},
	}

	mockRepo.On("List", ctx).Return(expectedAdditions, nil)

	result, err := service.ListAdditions(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedAdditions, result)
	mockRepo.AssertExpectations(t)
}
