package controller

import (
	config "github.com/Rifq11/Trava-be/config"
	"github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDashboardStatistics(c *gin.Context) {
	var totalDestinations int64
	var totalActiveOrders int64
	var totalRegisteredUsers int64

	config.DB.Table("destinations").Count(&totalDestinations)
	config.DB.Table("bookings").Where("status_id IN (1,2,5)").Count(&totalActiveOrders)
	config.DB.Table("users").Where("role_id = 2").Count(&totalRegisteredUsers)

	response := models.DashboardStatisticsResponse{
		TotalDestinations:   totalDestinations,
		TotalActiveOrders:   totalActiveOrders,
		TotalRegisteredUser: totalRegisteredUsers,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func GetMonthlySales(c *gin.Context) {
	destinationID := c.Query("destination_id")
	var results []models.MonthlySalesResponse

	query := config.DB.Table("bookings").
		Select("MONTH(created_at) AS month, SUM(total_price) AS revenue").
		Group("MONTH(created_at)").
		Order("MONTH(created_at) ASC")

	if destinationID != "" {
		query = query.Where("destination_id = ?", destinationID)
	}

	if err := query.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
