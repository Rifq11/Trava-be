package controller

import (
	"net/http"

	config "github.com/Rifq11/Trava-be/config"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func GetAllAccommodations(c *gin.Context) {
	var rows []models.AccommodationResponse

	err := config.DB.Table("destinations").
		Select(`
			destinations.id AS destination_id,
			destinations.name AS destination_name,
			destinations.image AS destination_image,
			transportation.id AS transport_id,
			transport_types.name AS transport_type_name,
			transportation.price,
			transportation.estimate,
			transportation.detail_transportation AS detail_transportation
		`).
		Joins("LEFT JOIN transportation ON destinations.id = transportation.destination_id").
		Joins("LEFT JOIN transport_types ON transportation.transport_type_id = transport_types.id").
		Order("destinations.id ASC").
		Scan(&rows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rows})
}

func CreateTransportation(c *gin.Context) {
	var req models.CreateTransportationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transport := models.Transportation{
		DestinationID:        req.DestinationID,
		TransportTypeID:      req.TransportTypeID,
		Price:                req.Price,
		Estimate:             req.Estimate,
		DetailTransportation: req.DetailTransportation,
	}

	if err := config.DB.Create(&transport).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transportation created successfully",
		"data":    transport,
	})
}

func UpdateTransportation(c *gin.Context) {
	id := c.Param("id")

	var transport models.Transportation
	if err := config.DB.First(&transport, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transportation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req models.UpdateTransportationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Estimate != nil {
		updates["estimate"] = *req.Estimate
	}
	if req.DetailTransportation != nil {
		updates["detail_transportation"] = *req.DetailTransportation
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := config.DB.Model(&transport).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config.DB.First(&transport, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Transportation updated successfully",
		"data":    transport,
	})
}

func DeleteTransportation(c *gin.Context) {
	id := c.Param("id")

	var transport models.Transportation
	if err := config.DB.First(&transport, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transportation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// delete reviews
	if err := config.DB.Where(
		"booking_id IN (SELECT id FROM bookings WHERE transportation_id = ?)", id,
	).Delete(&models.Review{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related reviews"})
		return
	}

	// delete payments
	if err := config.DB.Where(
		"booking_id IN (SELECT id FROM bookings WHERE transportation_id = ?)", id,
	).Delete(&models.Payment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related payments"})
		return
	}

	// delete bookings
	if err := config.DB.Where("transportation_id = ?", id).Delete(&models.Booking{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related bookings"})
		return
	}

	// delete transportation
	if err := config.DB.Delete(&transport).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transportation, bookings, reviews, and payments deleted successfully",
	})
}