package handlers

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/domain/category"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterCategoryRoutes(router *gin.RouterGroup, service services.CategoryService) {
	categories := router.Group("/categories")
	{
		categories.POST("/", CreateCategoryHandler(service))
		categories.GET("/:id", GetCategoryByIDHandler(service))
		categories.PUT("/:id", UpdateCategoryHandler(service))
		categories.DELETE("/:id", DeleteCategoryHandler(service))
		categories.GET("/", ListCategoriesHandler(service))
	}
}

// @Summary Create a Category
// @Description Cria uma nova categoria
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param category body category.Category true "Categoria a ser criada"
// @Success 201 {object} map[string]category.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /categories [post]
func CreateCategoryHandler(service services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cte category.Category
		if err := c.ShouldBindJSON(&cte); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := service.CreateCategory(c.Request.Context(), &cte)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, cte)
	}
}

// @Summary Get Category by ID
// @Description Recupera uma única categoria pelo seu ID
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path string true "ID da Categoria"
// @Success 200 {object} map[string]category.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /categories/{id} [get]
func GetCategoryByIDHandler(service services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
			return
		}

		cte, err := service.GetCategoryByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"category": cte})
	}
}

// @Summary Update a Category
// @Description Atualiza uma categoria existente pelo ID
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path string true "ID da Categoria"
// @Param category body category.Category true "Categoria a ser atualizada"
// @Success 200 {object} map[string]category.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /categories/{id} [put]
func UpdateCategoryHandler(service services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
			return
		}

		var cte category.Category
		if err := c.ShouldBindJSON(&cte); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cte.ID = id

		err = service.UpdateCategory(c.Request.Context(), &cte)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, cte)
	}
}

// @Summary Delete a Category
// @Description Deleta uma categoria pelo ID
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param id path string true "ID da Categoria"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /categories/{id} [delete]
func DeleteCategoryHandler(service services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID da categoria inválido"})
			return
		}

		err = service.DeleteCategory(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary List Categories
// @Description Recupera uma lista de todas as categorias
// @Tags Categories
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string][]category.Category
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /categories [get]
func ListCategoriesHandler(service services.CategoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := service.ListCategories(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}
