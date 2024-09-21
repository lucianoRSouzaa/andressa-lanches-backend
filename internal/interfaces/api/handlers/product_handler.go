package handlers

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/domain/product"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterProductRoutes(router *gin.RouterGroup, service services.ProductService) {
	products := router.Group("/products")
	{
		products.POST("/", CreateProductHandler(service))
		products.GET("/:id", GetProductByIDHandler(service))
		products.PUT("/:id", UpdateProductHandler(service))
		products.DELETE("/:id", DeleteProductHandler(service))
		products.GET("/", ListProductsHandler(service))
	}
}

// @Summary Create a Product
// @Description Cria um novo produto
// @Tags Products
// @Accept  json
// @Produce  json
// @Param product body product.Product true "Produto a ser criado"
// @Success 201 {object} map[string]product.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /products [post]
func CreateProductHandler(service services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p product.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := service.CreateProduct(c.Request.Context(), &p)
		if err != nil {
			switch err {
			case product.ErrProductNameRequired, product.ErrProductPricePositive, product.ErrProductCategoryID:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, p)
	}
}

// @Summary Get Product by ID
// @Description Recupera um único produto pelo seu ID
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path string true "ID do Produto"
// @Success 200 {object} map[string]product.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /products/{id} [get]
func GetProductByIDHandler(service services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
			return
		}

		product, err := service.GetProductByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"product": product})
	}
}

// @Summary Update a Product
// @Description Atualiza um produto existente pelo ID
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path string true "ID do Produto"
// @Param product body product.Product true "Produto a ser atualizado"
// @Success 200 {object} map[string]product.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /products/{id} [put]
func UpdateProductHandler(service services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
			return
		}

		var p product.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		p.ID = id

		err = service.UpdateProduct(c.Request.Context(), &p)
		if err != nil {
			switch err {
			case product.ErrProductNameRequired, product.ErrProductPricePositive, product.ErrProductCategoryID:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case product.ErrProductNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, p)
	}
}

// @Summary Delete a Product
// @Description Deleta um produto pelo ID
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path string true "ID do Produto"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /products/{id} [delete]
func DeleteProductHandler(service services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do produto inválido"})
			return
		}

		err = service.DeleteProduct(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary List Products
// @Description Recupera uma lista de todos os produtos
// @Tags Products
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string][]product.Product
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /products [get]
func ListProductsHandler(service services.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := service.ListProducts(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}
