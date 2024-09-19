package handlers

import (
	"andressa-lanches/internal/application/services"
	"andressa-lanches/internal/domain/addition"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterAdditionRoutes(router *gin.RouterGroup, service services.AdditionService) {
	additions := router.Group("/additions")
	{
		additions.POST("/", CreateAdditionHandler(service))
		additions.GET("/:id", GetAdditionByIDHandler(service))
		additions.PUT("/:id", UpdateAdditionHandler(service))
		additions.DELETE("/:id", DeleteAdditionHandler(service))
		additions.GET("/", ListAdditionsHandler(service))
	}
}

// @Summary Create an Addition
// @Description Cria um novo acréscimo
// @Tags Additions
// @Accept  json
// @Produce  json
// @Param addition body addition.Addition true "Acréscimo a ser criado"
// @Success 201 {object} map[string]addition.Addition
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /additions [post]
func CreateAdditionHandler(service services.AdditionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var add addition.Addition
		if err := c.ShouldBindJSON(&add); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := service.CreateAddition(c.Request.Context(), &add)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, add)
	}
}

// @Summary Get Addition by ID
// @Description Recupera um único acréscimo pelo seu ID
// @Tags Additions
// @Accept  json
// @Produce  json
// @Param id path string true "ID do Acréscimo"
// @Success 200 {object} map[string]addition.Addition
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /additions/{id} [get]
func GetAdditionByIDHandler(service services.AdditionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do acréscimo inválido"})
			return
		}

		add, err := service.GetAdditionByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"addition": add})
	}
}

// @Summary Update an Addition
// @Description Atualiza um acréscimo existente pelo ID
// @Tags Additions
// @Accept  json
// @Produce  json
// @Param id path string true "ID do Acréscimo"
// @Param addition body addition.Addition true "Acréscimo a ser atualizado"
// @Success 200 {object} map[string]addition.Addition
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /additions/{id} [put]
func UpdateAdditionHandler(service services.AdditionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do acréscimo inválido"})
			return
		}

		var add addition.Addition
		if err := c.ShouldBindJSON(&add); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		add.ID = id

		err = service.UpdateAddition(c.Request.Context(), &add)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, add)
	}
}

// @Summary Delete an Addition
// @Description Deleta um acréscimo pelo ID
// @Tags Additions
// @Accept  json
// @Produce  json
// @Param id path string true "ID do Acréscimo"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /additions/{id} [delete]
func DeleteAdditionHandler(service services.AdditionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID do acréscimo inválido"})
			return
		}

		err = service.DeleteAddition(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary List Additions
// @Description Recupera uma lista de todos os acréscimos
// @Tags Additions
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string][]addition.Addition
// @Failure 500 {object} map[string]string
// @Router /additions [get]
func ListAdditionsHandler(service services.AdditionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		additions, err := service.ListAdditions(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"additions": additions})
	}
}
