package tests

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/config"
	"andressa-lanches/internal/domain/addition"
	"andressa-lanches/internal/domain/product"
	"andressa-lanches/internal/domain/sale"
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

func setupSaleTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	config.JWTSecret = "test_secret"
	config.AuthUser = "test_user"
	config.AuthPassword = "test_password"

	// Repositórios em memória
	saleRepo := repository.NewInMemorySaleRepository()
	productRepo := repository.NewInMemoryProductRepository()
	additionRepo := repository.NewInMemoryAdditionRepository()

	// Serviços
	saleService := services.NewSaleService(saleRepo, productRepo, additionRepo)
	productService := services.NewProductService(productRepo)
	additionService := services.NewAdditionService(additionRepo)

	router := gin.Default()
	router.POST("/auth/login", handlers.LoginHandler())

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())

	// Registrar rotas necessárias
	handlers.RegisterSaleRoutes(protected, saleService)
	handlers.RegisterProductRoutes(protected, productService)
	handlers.RegisterAdditionRoutes(protected, additionService)

	return router
}

func TestCreateSale_Success(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Criar um produto
	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       10.0,
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

	// Criar um acréscimo
	newAddition := &addition.Addition{
		Name:  "Extra Cheese",
		Price: 2.5,
	}
	payload, _ = json.Marshal(newAddition)

	req, _ = http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdAddition addition.Addition
	err = json.Unmarshal(w.Body.Bytes(), &createdAddition)
	require.NoError(t, err)

	// Criar uma venda
	newSale := &sale.Sale{
		Discount:          0.0,
		AdditionalCharges: 0.0,
		Items: []sale.SaleItem{
			{
				ProductID: createdProduct.ID,
				Quantity:  2,
				Additions: []addition.Addition{
					{ID: createdAddition.ID},
				},
			},
			{
				ProductID: createdProduct.ID,
				Quantity:  1,
				Additions: []addition.Addition{
					{ID: createdAddition.ID},
					{ID: createdAddition.ID},
				},
			},
		},
	}
	payload, _ = json.Marshal(newSale)

	req, _ = http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdSale sale.Sale
	err = json.Unmarshal(w.Body.Bytes(), &createdSale)
	assert.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, createdSale.ID)
	assert.Equal(t, 1, len(createdSale.Items))
	assert.Equal(t, createdProduct.ID, createdSale.Items[0].ProductID)
	assert.Equal(t, 2, createdSale.Items[0].Quantity)
	assert.Equal(t, createdProduct.Price, createdSale.Items[0].UnitPrice)
	assert.Equal(t, (createdProduct.Price+createdAddition.Price)*3, createdSale.Items[0].TotalPrice)
	assert.Equal(t, createdSale.TotalAmount, createdSale.Items[0].TotalPrice)
}

func TestCreateSale_ProductNotFound(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Criar uma venda com produto inexistente
	newSale := &sale.Sale{
		Items: []sale.SaleItem{
			{
				ProductID: uuid.New(),
				Quantity:  1,
			},
		},
	}
	payload, _ := json.Marshal(newSale)

	req, _ := http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "produto não encontrado")
}

func TestCreateSale_AdditionNotFound(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Criar um produto
	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       10.0,
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

	// Criar uma venda com acréscimo inexistente
	newSale := &sale.Sale{
		Items: []sale.SaleItem{
			{
				ProductID: createdProduct.ID,
				Quantity:  1,
				Additions: []addition.Addition{
					{ID: uuid.New()},
				},
			},
		},
	}
	payload, _ = json.Marshal(newSale)

	req, _ = http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "acréscimo não encontrado")
}

func TestCreateSale_InvalidQuantity(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Criar um produto
	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       10.0,
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

	// Criar uma venda com quantidade inválida
	newSale := &sale.Sale{
		Items: []sale.SaleItem{
			{
				ProductID: createdProduct.ID,
				Quantity:  0,
			},
		},
	}
	payload, _ = json.Marshal(newSale)

	req, _ = http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "a quantidade deve ser positiva")
}

func TestGetSaleByID_Success(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Criar um produto
	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       10.0,
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

	// Criar um acréscimo
	newAddition := &addition.Addition{
		Name:  "Extra Cheese",
		Price: 2.5,
	}
	payload, _ = json.Marshal(newAddition)

	req, _ = http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdAddition addition.Addition
	err = json.Unmarshal(w.Body.Bytes(), &createdAddition)
	require.NoError(t, err)

	// Criar uma venda
	newSale := &sale.Sale{
		Discount:          0.0,
		AdditionalCharges: 0.0,
		Items: []sale.SaleItem{
			{
				ProductID: createdProduct.ID,
				Quantity:  2,
				Additions: []addition.Addition{
					{ID: createdAddition.ID},
				},
			},
		},
	}

	payload, _ = json.Marshal(newSale)

	req, _ = http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdSale sale.Sale
	err = json.Unmarshal(w.Body.Bytes(), &createdSale)
	assert.NoError(t, err)

	// Buscar a venda pelo ID
	req, _ = http.NewRequest(http.MethodGet, "/sales/"+createdSale.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse map[string]sale.Sale
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.NoError(t, err)
	fetchedSale := getResponse["sale"]

	assert.Equal(t, createdSale.ID.String(), fetchedSale.ID.String())
}

func TestGetSaleByID_NotFound(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Buscar uma venda inexistente
	req, _ := http.NewRequest(http.MethodGet, "/sales/"+uuid.New().String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "venda não encontrada")
}

func TestDeleteSale_Success(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Criar um produto
	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       10.0,
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

	// Criar um acréscimo
	newAddition := &addition.Addition{
		Name:  "Extra Cheese",
		Price: 2.5,
	}
	payload, _ = json.Marshal(newAddition)

	req, _ = http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdAddition addition.Addition
	err = json.Unmarshal(w.Body.Bytes(), &createdAddition)
	require.NoError(t, err)

	// Criar uma venda
	newSale := &sale.Sale{
		Discount:          0.0,
		AdditionalCharges: 0.0,
		Items: []sale.SaleItem{
			{
				ProductID: createdProduct.ID,
				Quantity:  2,
				Additions: []addition.Addition{
					{ID: createdAddition.ID},
				},
			},
		},
	}

	payload, _ = json.Marshal(newSale)

	req, _ = http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdSale sale.Sale
	err = json.Unmarshal(w.Body.Bytes(), &createdSale)
	assert.NoError(t, err)

	// Deletar a venda
	req, _ = http.NewRequest(http.MethodDelete, "/sales/"+createdSale.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Tentar buscar a venda deletada
	req, _ = http.NewRequest(http.MethodGet, "/sales/"+createdSale.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListSales_Success(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	expectedNumberOfSales := 3

	// Criar um produto
	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       10.0,
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

	// Criar um acréscimo
	newAddition := &addition.Addition{
		Name:  "Extra Cheese",
		Price: 2.5,
	}
	payload, _ = json.Marshal(newAddition)

	req, _ = http.NewRequest(http.MethodPost, "/additions/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdAddition addition.Addition
	err = json.Unmarshal(w.Body.Bytes(), &createdAddition)
	require.NoError(t, err)

	// Criar algumas vendas
	for i := 0; i <= expectedNumberOfSales; i++ {
		newSale := &sale.Sale{
			Discount:          0.0,
			AdditionalCharges: 0.0,
			Items: []sale.SaleItem{
				{
					ProductID: createdProduct.ID,
					Quantity:  2,
					Additions: []addition.Addition{
						{ID: createdAddition.ID},
					},
				},
			},
		}

		payload, _ = json.Marshal(newSale)

		req, _ = http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	}

	// Listar as vendas
	req, _ = http.NewRequest(http.MethodGet, "/sales/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse map[string][]sale.Sale
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.NoError(t, err)
	fetchedSales := getResponse["sales"]

	// Verifique se o número de vendas corresponde
	assert.GreaterOrEqual(t, len(fetchedSales), expectedNumberOfSales)
}

func TestCreateSale_InvalidProductID(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Criar uma venda com ID de produto inválido
	newSale := &sale.Sale{
		Items: []sale.SaleItem{
			{
				ProductID: uuid.Nil,
				Quantity:  1,
			},
		},
	}
	payload, _ := json.Marshal(newSale)

	req, _ := http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "ID do produto é obrigatório para o item da venda")
}

func TestCreateSale_InvalidAdditionID(t *testing.T) {
	router := setupSaleTestRouter()
	token := getValidToken(t, router)

	// Criar um produto
	newProduct := &product.Product{
		Name:        "Test Product",
		Description: "A product for testing",
		Price:       10.0,
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

	// Criar uma venda com ID de acréscimo inválido
	newSale := &sale.Sale{
		Items: []sale.SaleItem{
			{
				ProductID: createdProduct.ID,
				Quantity:  1,
				Additions: []addition.Addition{
					{ID: uuid.Nil},
				},
			},
		},
	}
	payload, _ = json.Marshal(newSale)

	req, _ = http.NewRequest(http.MethodPost, "/sales/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "ID do acréscimo inválido")
}
