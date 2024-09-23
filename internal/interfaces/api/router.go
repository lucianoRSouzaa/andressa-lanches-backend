package api

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/interfaces/api/docs"
	"andressa-lanches/internal/interfaces/api/handlers"
	"andressa-lanches/internal/interfaces/api/middlewares"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	productService services.ProductService,
	categoryService services.CategoryService,
	additionService services.AdditionService,
	saleService services.SaleService,
) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggingMiddleware())

	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	router.Use(p.Instrument())

	auth := router.Group("/auth")
	{
		auth.POST("/login", handlers.LoginHandler())
	}

	// Rotas com JWT
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Produtos
		handlers.RegisterProductRoutes(protected, productService)

		// Categorias
		handlers.RegisterCategoryRoutes(protected, categoryService)

		// Acréscimos
		handlers.RegisterAdditionRoutes(protected, additionService)

		// Vendas
		handlers.RegisterSaleRoutes(protected, saleService)
	}

	docs.InitializeSwagger(router)

	return router
}
