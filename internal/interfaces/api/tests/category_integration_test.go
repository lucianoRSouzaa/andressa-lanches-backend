package tests

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/config"
	"andressa-lanches/internal/domain/category"
	"andressa-lanches/internal/infrastructure/repository"
	"andressa-lanches/internal/interfaces/api/handlers"
	"andressa-lanches/internal/interfaces/api/middlewares"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCategoryTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	config.JWTSecret = "test_secret"
	config.AuthUser = "test_user"
	config.AuthPassword = "test_password"

	categoryRepo := repository.NewInMemoryCategoryRepository()
	categoryService := services.NewCategoryService(categoryRepo)

	router := gin.Default()
	router.POST("/auth/login", handlers.LoginHandler())

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	handlers.RegisterCategoryRoutes(protected, categoryService)

	return router
}

func TestCreateCategory_Success(t *testing.T) {
	router := setupCategoryTestRouter()
	token := getValidToken(t, router)

	newCategory := &category.Category{
		Name:        "Test Category",
		Description: "A category for testing",
	}
	payload, _ := json.Marshal(newCategory)

	req, _ := http.NewRequest(http.MethodPost, "/categories/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdCategory category.Category
	err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
	assert.NoError(t, err)

	assert.Equal(t, newCategory.Name, createdCategory.Name)
	assert.Equal(t, newCategory.Description, createdCategory.Description)
	assert.NotEqual(t, uuid.Nil, createdCategory.ID)
}

func TestCreateCategory_WithoutNameError(t *testing.T) {
	router := setupCategoryTestRouter()
	token := getValidToken(t, router)

	newCategory := &category.Category{
		Name:        "",
		Description: "A category for testing",
	}
	payload, _ := json.Marshal(newCategory)

	req, _ := http.NewRequest(http.MethodPost, "/categories/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, category.ErrCategoryNameRequired.Error(), errorResponse.Error)
}

func TestGetCategoryByID_Success(t *testing.T) {
	router := setupCategoryTestRouter()
	token := getValidToken(t, router)

	// Criar uma categoria primeiro
	newCategory := &category.Category{
		Name:        "Test Category",
		Description: "A category for testing",
	}
	payload, _ := json.Marshal(newCategory)

	req, _ := http.NewRequest(http.MethodPost, "/categories/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdCategory category.Category
	err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
	require.NoError(t, err)

	// Buscar a categoria pelo ID
	req, _ = http.NewRequest(http.MethodGet, "/categories/"+createdCategory.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse map[string]category.Category
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.NoError(t, err)
	fetchedCategory := getResponse["category"]

	assert.Equal(t, createdCategory.ID, fetchedCategory.ID)
	assert.Equal(t, createdCategory.Name, fetchedCategory.Name)
	assert.Equal(t, createdCategory.Description, fetchedCategory.Description)
}

func TestUpdateCategory_Success(t *testing.T) {
	router := setupCategoryTestRouter()
	token := getValidToken(t, router)

	// Criar uma categoria
	newCategory := &category.Category{
		Name:        "Old Name",
		Description: "Old Description",
	}
	payload, _ := json.Marshal(newCategory)

	req, _ := http.NewRequest(http.MethodPost, "/categories/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdCategory category.Category
	err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
	require.NoError(t, err)

	// Atualizar a categoria
	updatedCategory := &category.Category{
		Name:        "New Name",
		Description: "New Description",
	}
	payload, _ = json.Marshal(updatedCategory)

	req, _ = http.NewRequest(http.MethodPut, "/categories/"+createdCategory.ID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedCategoryResponse category.Category
	err = json.Unmarshal(w.Body.Bytes(), &updatedCategoryResponse)
	assert.NoError(t, err)

	assert.Equal(t, createdCategory.ID, updatedCategoryResponse.ID)
	assert.Equal(t, updatedCategory.Name, updatedCategoryResponse.Name)
	assert.Equal(t, updatedCategory.Description, updatedCategoryResponse.Description)
}

func TestDeleteCategory_Success(t *testing.T) {
	router := setupCategoryTestRouter()
	token := getValidToken(t, router)

	// Criar uma categoria
	newCategory := &category.Category{
		Name:        "To be deleted",
		Description: "This category will be deleted",
	}
	payload, _ := json.Marshal(newCategory)

	req, _ := http.NewRequest(http.MethodPost, "/categories/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdCategory category.Category
	err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
	require.NoError(t, err)

	// Deletar a categoria
	req, _ = http.NewRequest(http.MethodDelete, "/categories/"+createdCategory.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Tentar buscar a categoria deletada
	req, _ = http.NewRequest(http.MethodGet, "/categories/"+createdCategory.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListCategories_Success(t *testing.T) {
	router := setupCategoryTestRouter()
	token := getValidToken(t, router)

	// Criar algumas categorias
	categories := []category.Category{
		{Name: "Category 1", Description: "First category"},
		{Name: "Category 2", Description: "Second category"},
	}

	for _, c := range categories {
		payload, _ := json.Marshal(c)
		req, _ := http.NewRequest(http.MethodPost, "/categories/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Listar as categorias
	req, _ := http.NewRequest(http.MethodGet, "/categories/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var listResponse map[string][]category.Category
	err := json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.NoError(t, err)
	fetchedCategories := listResponse["categories"]

	assert.Len(t, fetchedCategories, len(categories))
}
