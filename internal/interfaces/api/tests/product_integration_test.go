package tests

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/config"
	"andressa-lanches/internal/domain/product"
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

type ErrorResponse struct {
	Error string `json:"error"`
}

func setupProductTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	config.JWTSecret = "test_secret"
	config.AuthUser = "test_user"
	config.AuthPassword = "test_password"

	productRepo := repository.NewInMemoryProductRepository()
	productService := services.NewProductService(productRepo)

	router := gin.Default()

	router.POST("/auth/login", handlers.LoginHandler())

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	handlers.RegisterProductRoutes(protected, productService)

	return router
}

func getValidToken(t *testing.T, router *gin.Engine) string {
	loginData := map[string]string{
		"username": "test_user",
		"password": "test_password",
	}
	payload, _ := json.Marshal(loginData)

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotEmpty(t, response["token"])

	return response["token"]
}

func TestCreateProduct_Success(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       9.99,
		CategoryID:  uuid.New(),
	}
	payload, _ := json.Marshal(newProduct)

	req, _ := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdProduct product.Product
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	assert.NoError(t, err)

	assert.Equal(t, newProduct.Name, createdProduct.Name)
	assert.Equal(t, newProduct.Description, createdProduct.Description)
	assert.Equal(t, newProduct.Price, createdProduct.Price)
	assert.NotEqual(t, uuid.Nil, createdProduct.ID)
}

func TestCreateProduct_InvalidInput(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	// Dados inválidos (faltando campos obrigatórios)
	invalidProduct := map[string]interface{}{
		"Price": 9.99,
	}
	payload, _ := json.Marshal(invalidProduct)

	req, _ := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProduct_WithoutCategory(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	// Dados inválidos (faltando campo category_id)
	invalidProduct := map[string]interface{}{
		"Price":       9.99,
		"Name":        "Test Product",
		"Description": "A product for testing",
	}

	payload, _ := json.Marshal(invalidProduct)

	req, _ := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "o ID da categoria do produto é obrigatório", errorResponse.Error)
}

func TestGetProductByID_Success(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	// Criar um produto primeiro
	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       9.99,
		CategoryID:  uuid.New(),
	}
	payload, _ := json.Marshal(newProduct)

	req, _ := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdProduct product.Product
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	require.NoError(t, err)

	// Buscar o produto pelo ID
	req, _ = http.NewRequest(http.MethodGet, "/products/"+createdProduct.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse map[string]product.Product
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.NoError(t, err)
	fetchedProduct := getResponse["product"]

	assert.Equal(t, createdProduct.ID, fetchedProduct.ID)
	assert.Equal(t, createdProduct.Name, fetchedProduct.Name)
	assert.Equal(t, createdProduct.Description, fetchedProduct.Description)
	assert.Equal(t, createdProduct.Price, fetchedProduct.Price)
}

func TestGetProductByID_NotFound(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	// Buscar um produto com ID inexistente
	nonExistentID := uuid.New().String()
	req, _ := http.NewRequest(http.MethodGet, "/products/"+nonExistentID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateProduct_Success(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	// Criar um produto
	newProduct := &product.Product{
		Name:        "Old Name",
		Description: "Old Description",
		Price:       9.99,
		CategoryID:  uuid.New(),
	}
	payload, _ := json.Marshal(newProduct)

	req, _ := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdProduct product.Product
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	require.NoError(t, err)

	// Atualizar o produto
	updatedProduct := &product.Product{
		Name:        "New Name",
		Description: "New Description",
		Price:       19.99,
		CategoryID:  createdProduct.CategoryID,
	}
	payload, _ = json.Marshal(updatedProduct)

	req, _ = http.NewRequest(http.MethodPut, "/products/"+createdProduct.ID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedProductResponse product.Product
	err = json.Unmarshal(w.Body.Bytes(), &updatedProductResponse)
	assert.NoError(t, err)

	assert.Equal(t, createdProduct.ID, updatedProductResponse.ID)
	assert.Equal(t, updatedProduct.Name, updatedProductResponse.Name)
	assert.Equal(t, updatedProduct.Description, updatedProductResponse.Description)
	assert.Equal(t, updatedProduct.Price, updatedProductResponse.Price)
	assert.Equal(t, createdProduct.CategoryID, updatedProductResponse.CategoryID)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	// Atualizar um produto inexistente
	updatedProduct := &product.Product{
		Name:        "New Name",
		Description: "New Description",
		Price:       19.99,
	}
	payload, _ := json.Marshal(updatedProduct)

	nonExistentID := uuid.New().String()
	req, _ := http.NewRequest(http.MethodPut, "/products/"+nonExistentID, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteProduct_Success(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	// Criar um produto
	newProduct := &product.Product{
		Name:        "To be deleted",
		Description: "This product will be deleted",
		Price:       9.99,
		CategoryID:  uuid.New(),
	}
	payload, _ := json.Marshal(newProduct)

	req, _ := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdProduct product.Product
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	require.NoError(t, err)

	// Deletar o produto
	req, _ = http.NewRequest(http.MethodDelete, "/products/"+createdProduct.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Tentar buscar o produto deletado
	req, _ = http.NewRequest(http.MethodGet, "/products/"+createdProduct.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListProducts_Success(t *testing.T) {
	router := setupProductTestRouter()
	token := getValidToken(t, router)

	// Criar alguns produtos
	products := []product.Product{
		{Name: "Product 1", Description: "First product", Price: 10.0, CategoryID: uuid.New()},
		{Name: "Product 2", Description: "Second product", Price: 20.0, CategoryID: uuid.New()},
	}

	for _, p := range products {
		payload, _ := json.Marshal(p)
		req, _ := http.NewRequest(http.MethodPost, "/products/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Listar os produtos
	req, _ := http.NewRequest(http.MethodGet, "/products/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var listResponse map[string][]product.Product
	err := json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.NoError(t, err)
	fetchedProducts := listResponse["products"]

	assert.Len(t, fetchedProducts, len(products))
}

func TestProtectedRoutes_Unauthorized(t *testing.T) {
	router := setupProductTestRouter()

	// Tentar acessar uma rota protegida sem token
	req, _ := http.NewRequest(http.MethodGet, "/products/", nil)
	// Não definir o header Authorization

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
