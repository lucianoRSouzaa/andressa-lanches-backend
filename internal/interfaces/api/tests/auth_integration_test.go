package tests

import (
	"andressa-lanches/internal/config"
	"andressa-lanches/internal/interfaces/api/handlers"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.JWTSecret = "test_secret"
	config.AuthUser = "test_user"
	config.AuthPassword = "test_password"

	// Configurar o roteador
	router := gin.Default()
	router.POST("/auth/login", handlers.LoginHandler())

	// Dados de login válidos
	loginData := map[string]string{
		"username": "test_user",
		"password": "test_password",
	}
	payload, _ := json.Marshal(loginData)

	// Criar uma requisição HTTP
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Executar a requisição
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verificar a resposta
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.JWTSecret = "test_secret"
	router := gin.Default()
	router.POST("/auth/login", handlers.LoginHandler())

	// Dados de login inválidos
	loginData := map[string]string{
		"username": "admin",
		"password": "wrong_password",
	}
	payload, _ := json.Marshal(loginData)

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
