package controller

import (
	"net/http"
	"strings"

	config "github.com/Rifq11/Trava-be/config"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdInt := userID.(int)

	var destination models.Destination
	result := config.DB.First(&destination, req.DestinationID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destination or transportation not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	var transportation models.Transportation
	result = config.DB.First(&transportation, req.TransportationID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destination or transportation not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	destinationPrice := destination.PricePerPerson * req.PeopleCount
	transportPrice := transportation.Price
	totalPrice := destinationPrice + transportPrice

	booking := models.Booking{
		UserID:           userIdInt,
		DestinationID:    req.DestinationID,
		TransportationID: req.TransportationID,
		PaymentMethodID:  req.PaymentMethodID,
		StatusID:         1, // Pending
		PeopleCount:      req.PeopleCount,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		TransportPrice:   transportPrice,
		DestinationPrice: destinationPrice,
		TotalPrice:       totalPrice,
	}

	result = config.DB.Create(&booking)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking created successfully",
		"data":    booking,
	})
}

func GetMyBookings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdInt := userID.(int)

	var bookings []models.Booking
	result := config.DB.
		Where("user_id = ?", userIdInt).
		Find(&bookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if bookings == nil {
		bookings = []models.Booking{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
	})
}

func GetAllBookingsAdmin(c *gin.Context) {
	var bookings []models.Booking
	if err := config.DB.
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if bookings == nil {
		bookings = []models.Booking{}
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

func ApproveBooking(c *gin.Context) {
	updateBookingStatus(c, "approved")
}

func RejectBooking(c *gin.Context) {
	updateBookingStatus(c, "rejected")
}

func updateBookingStatus(c *gin.Context, statusName string) {
	id := c.Param("id")

	var booking models.Booking
	if err := config.DB.First(&booking, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var status models.BookingStatus
	if err := config.DB.Table("booking_status").
		Where("LOWER(name) = ?", strings.ToLower(statusName)).
		First(&status).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking status"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&booking).Update("status_id", status.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config.DB.First(&booking, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking status updated",
		"data":    booking,
	})
}
