package controller

import (
	"net/http"

	config "github.com/Rifq11/Trava-be/config"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
)

func GetPaymentMethods(c *gin.Context) {
	var paymentMethods []models.PaymentMethod
	result := config.DB.Find(&paymentMethods)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": paymentMethods,
	})
}
