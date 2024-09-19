package handlers

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/domain/sale"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterSaleRoutes(router *gin.RouterGroup, service services.SaleService) {
	sales := router.Group("/sales")
	{
		sales.POST("/", CreateSaleHandler(service))
		sales.GET("/:id", GetSaleByIDHandler(service))
		sales.GET("/", ListSalesHandler(service))
		sales.DELETE("/:id", DeleteSaleHandler(service))
	}
}

// @Summary Create a Sale
// @Description Cria uma nova venda
// @Tags Sales
// @Accept  json
// @Produce  json
// @Param sale body sale.Sale true "Venda a ser criada"
// @Success 201 {object} map[string]sale.Sale
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /sales [post]
func CreateSaleHandler(service services.SaleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s sale.Sale
		if err := c.ShouldBindJSON(&s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := service.CreateSale(c.Request.Context(), &s)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, s)
	}
}

// @Summary Get Sale by ID
// @Description Recupera uma única venda pelo seu ID
// @Tags Sales
// @Accept  json
// @Produce  json
// @Param id path string true "ID da Venda"
// @Success 200 {object} map[string]sale.Sale
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /sales/{id} [get]
func GetSaleByIDHandler(service services.SaleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID da venda inválido"})
			return
		}

		s, err := service.GetSaleByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sale": s})
	}
}

// @Summary Delete a Sale
// @Description Deleta uma venda pelo ID
// @Tags Sales
// @Accept  json
// @Produce  json
// @Param id path string true "ID da Venda"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /sales/{id} [delete]
func DeleteSaleHandler(service services.SaleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID da venda inválido"})
			return
		}

		err = service.DeleteSale(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary List Sales
// @Description Recupera uma lista de todas as vendas
// @Tags Sales
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string][]sale.Sale
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /sales [get]
func ListSalesHandler(service services.SaleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sales, err := service.ListSales(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sales": sales})
	}
}
