package controller

import (
	"net/http"

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

	var bookings []models.BookingResponse
	result := config.DB.
		Table("bookings").
		Select("bookings.id as booking_id, destinations.name as destination_name, destinations.location, bookings.people_count, bookings.start_date, bookings.end_date, bookings.total_price, booking_status.name as status_name, payment_methods.name as payment_method_name").
		Joins("INNER JOIN destinations ON bookings.destination_id = destinations.id").
		Joins("INNER JOIN booking_status ON bookings.status_id = booking_status.id").
		Joins("INNER JOIN payment_methods ON bookings.payment_method_id = payment_methods.id").
		Where("bookings.user_id = ?", userIdInt).
		Order("bookings.id DESC").
		Scan(&bookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
	})
}
