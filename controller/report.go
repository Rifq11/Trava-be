package controller

import (
	"strings"

	config "github.com/Rifq11/Trava-be/config"
	"github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
)

func GetReportOrders(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")

	var results []models.ReportOrderResponse

	query := config.DB.Table("bookings").
		Select(`
			bookings.id AS id,
			bookings.user_id,
			users.full_name AS user_name,
			destinations.name AS destination_name,
			bookings.start_date,
			bookings.end_date,
			bookings.people_count,
			bookings.transport_price,
			bookings.total_price,
			booking_status.name AS status_name
		`).
		Joins("JOIN users ON users.id = bookings.user_id").
		Joins("JOIN destinations ON destinations.id = bookings.destination_id").
		Joins("JOIN booking_status ON booking_status.id = bookings.status_id")

	if status != "" {
		query = query.Where("LOWER(booking_status.name) = ?", strings.ToLower(status))
	}

	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where(`
			users.full_name LIKE ? OR 
			destinations.name LIKE ? OR 
			CAST(bookings.id AS CHAR) LIKE ?`,
			pattern, pattern, pattern)
	}

	if err := query.Order("bookings.id DESC").Scan(&results).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": results})
}

func GetIncomeReport(c *gin.Context) {
	var income models.IncomeReportResponse

	err := config.DB.Table("bookings").
		Where("status_id = 5").
		Select("SUM(total_price) AS total_income").
		Scan(&income).Error

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": income})
}

func GetIncomeByDestination(c *gin.Context) {
	type IncomeDestination struct {
		DestinationName string  `json:"destination_name"`
		TotalIncome     float64 `json:"total_income"`
	}

	var results []IncomeDestination

	err := config.DB.Table("bookings").
		Select(`
			destinations.name AS destination_name,
			SUM(bookings.total_price) AS total_income
		`).
		Joins("JOIN destinations ON destinations.id = bookings.destination_id").
		Where("bookings.status_id = 5").
		Group("destinations.name").
		Order("total_income DESC").
		Scan(&results).Error

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": results})
}
