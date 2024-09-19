package main

// @title Andressa Lanches API
// @version 1.0
// @description API para gerenciar o sistema de lanches da Andressa.

// @contact.name Luciano Rodrigues De Souza
// @contact.email lucianorodriguess101@gmail.com

// @host localhost:3333
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Insira o token JWT no formato: Bearer {token}

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/config"
	"andressa-lanches/internal/infrastructure/db"
	"andressa-lanches/internal/infrastructure/repository"
	"andressa-lanches/internal/interfaces/api"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Configuração (carregar variáveis de ambiente, etc.)
	cfg := config.LoadConfig()

	pool, err := db.NewPgxPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	defer pool.Close()

	categoryRepo := repository.NewCategoryRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	additionRepo := repository.NewAdditionRepository(pool)
	saleRepo := repository.NewSaleRepository(pool)

	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo)
	additionService := services.NewAdditionService(additionRepo)
	saleService := services.NewSaleService(saleRepo, productRepo, additionRepo)

	router := api.SetupRouter(productService, categoryService, additionService, saleService)

	go func() {
		if err := router.Run(cfg.ServerAddress); err != nil {
			log.Fatalf("Falha ao iniciar o servidor: %v", err)
		}
	}()
	log.Printf("Servidor rodando em %s", cfg.ServerAddress)

	// Esperar por sinal de interrupção para encerrar
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Encerrando servidor...")
}
