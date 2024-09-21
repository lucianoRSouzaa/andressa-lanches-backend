package tests

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/config"
	"andressa-lanches/internal/domain/addition"
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

func setupAdditionTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	config.JWTSecret = "test_secret"
	config.AuthUser = "test_user"
	config.AuthPassword = "test_password"

	additionRepo := repository.NewInMemoryAdditionRepository()
	additionService := services.NewAdditionService(additionRepo)

	router := gin.Default()
	router.POST("/auth/login", handlers.LoginHandler())

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	handlers.RegisterAdditionRoutes(protected, additionService)

	return router
}

func TestCreateAddition_Success(t *testing.T) {
	router := setupAdditionTestRouter()
	token := getValidToken(t, router)

	newAddition := &addition.Addition{
		Name:  "Extra Cheese",
		Price: 2.5,
	}
	payload, _ := json.Marshal(newAddition)

	req, _ := http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdAddition addition.Addition
	err := json.Unmarshal(w.Body.Bytes(), &createdAddition)
	assert.NoError(t, err)

	assert.Equal(t, newAddition.Name, createdAddition.Name)
	assert.Equal(t, newAddition.Price, createdAddition.Price)
	assert.NotEqual(t, uuid.Nil, createdAddition.ID)
}

func TestCreateAddition_WithoutNameError(t *testing.T) {
	router := setupAdditionTestRouter()
	token := getValidToken(t, router)

	newAddition := &addition.Addition{
		Price: 2.5,
	}
	payload, _ := json.Marshal(newAddition)

	req, _ := http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, addition.ErrAdditionNameRequired.Error(), errorResponse.Error)
}

func TestCreateAddition_WithWrongPriceError(t *testing.T) {
	router := setupAdditionTestRouter()
	token := getValidToken(t, router)

	newAddition := &addition.Addition{
		Name:  "Extra Cheese",
		Price: -2.5,
	}
	payload, _ := json.Marshal(newAddition)

	req, _ := http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, addition.ErrAdditionPriceRequired.Error(), errorResponse.Error)
}

func TestGetAdditionByID_Success(t *testing.T) {
	router := setupAdditionTestRouter()
	token := getValidToken(t, router)

	// Criar um acréscimo primeiro
	newAddition := &addition.Addition{
		Name:  "Extra Bacon",
		Price: 3.0,
	}
	payload, _ := json.Marshal(newAddition)

	req, _ := http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdAddition addition.Addition
	err := json.Unmarshal(w.Body.Bytes(), &createdAddition)
	require.NoError(t, err)

	// Buscar o acréscimo pelo ID
	req, _ = http.NewRequest(http.MethodGet, "/additions/"+createdAddition.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse map[string]addition.Addition
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.NoError(t, err)
	fetchedAddition := getResponse["addition"]

	assert.Equal(t, createdAddition.ID, fetchedAddition.ID)
	assert.Equal(t, createdAddition.Name, fetchedAddition.Name)
	assert.Equal(t, createdAddition.Price, fetchedAddition.Price)
}

func TestUpdateAddition_Success(t *testing.T) {
	router := setupAdditionTestRouter()
	token := getValidToken(t, router)

	// Criar um acréscimo
	newAddition := &addition.Addition{
		Name:  "Old Name",
		Price: 1.0,
	}
	payload, _ := json.Marshal(newAddition)

	req, _ := http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdAddition addition.Addition
	err := json.Unmarshal(w.Body.Bytes(), &createdAddition)
	require.NoError(t, err)

	// Atualizar o acréscimo
	updatedAddition := &addition.Addition{
		Name:  "New Name",
		Price: 2.0,
	}
	payload, _ = json.Marshal(updatedAddition)

	req, _ = http.NewRequest(http.MethodPut, "/additions/"+createdAddition.ID.String(), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedAdditionResponse addition.Addition
	err = json.Unmarshal(w.Body.Bytes(), &updatedAdditionResponse)
	assert.NoError(t, err)

	assert.Equal(t, createdAddition.ID, updatedAdditionResponse.ID)
	assert.Equal(t, updatedAddition.Name, updatedAdditionResponse.Name)
	assert.Equal(t, updatedAddition.Price, updatedAdditionResponse.Price)
}

func TestDeleteAddition_Success(t *testing.T) {
	router := setupAdditionTestRouter()
	token := getValidToken(t, router)

	// Criar um acréscimo
	newAddition := &addition.Addition{
		Name:  "To be deleted",
		Price: 1.5,
	}
	payload, _ := json.Marshal(newAddition)

	req, _ := http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdAddition addition.Addition
	err := json.Unmarshal(w.Body.Bytes(), &createdAddition)
	require.NoError(t, err)

	// Deletar o acréscimo
	req, _ = http.NewRequest(http.MethodDelete, "/additions/"+createdAddition.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Tentar buscar o acréscimo deletado
	req, _ = http.NewRequest(http.MethodGet, "/additions/"+createdAddition.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListAdditions_Success(t *testing.T) {
	router := setupAdditionTestRouter()
	token := getValidToken(t, router)

	// Criar alguns acréscimos
	additions := []addition.Addition{
		{Name: "Addition 1", Price: 1.0},
		{Name: "Addition 2", Price: 2.0},
	}

	for _, a := range additions {
		payload, _ := json.Marshal(a)
		req, _ := http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Listar os acréscimos
	req, _ := http.NewRequest(http.MethodGet, "/additions/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var listResponse map[string][]addition.Addition
	err := json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.NoError(t, err)
	fetchedAdditions := listResponse["additions"]

	assert.Len(t, fetchedAdditions, len(additions))
}
