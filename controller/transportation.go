package controller

import (
	"net/http"

	config "github.com/Rifq11/Trava-be/config"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
)

func GetTransportationsByDestination(c *gin.Context) {
	destinationID := c.Param("id")

	var transportations []models.Transportation
	result := config.DB.Where("destination_id = ?", destinationID).Find(&transportations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transportations,
	})
}
