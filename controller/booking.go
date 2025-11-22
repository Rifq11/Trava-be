package controller

import (
	"net/http"
	"strings"
	"time"

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

	now := time.Now()
	var bookingsToUpdate []models.Booking
	config.DB.Where("user_id = ?", userIdInt).Find(&bookingsToUpdate)

	// completed status
	var completedStatus models.BookingStatus
	config.DB.Table("booking_status").
		Where("LOWER(name) = ?", "completed").
		First(&completedStatus)

	for _, booking := range bookingsToUpdate {
		var endDate time.Time
		var err error

		endDate, err = time.Parse("2006-01-02 15:04:05", booking.EndDate)
		if err != nil {
			endDate, err = time.Parse(time.RFC3339, booking.EndDate)
			if err != nil {
				endDate, err = time.Parse("2006-01-02", booking.EndDate)
				if err != nil {
					continue
				}
			}
		}

		if endDate.Before(now) || endDate.Equal(now) {
			var currentStatus models.BookingStatus
			if err := config.DB.Table("booking_status").Where("id = ?", booking.StatusID).First(&currentStatus).Error; err == nil {
				statusName := strings.ToLower(currentStatus.Name)
				if statusName == "pending" || statusName == "approved" {
					if completedStatus.ID > 0 {
						config.DB.Model(&booking).Update("status_id", completedStatus.ID)
					}
				}
			}
		}
	}

	var bookings []models.BookingResponse
	result := config.DB.Table("bookings"). // join to get other column in other table
						Select(`
			bookings.id as booking_id,
			bookings.destination_id as destination_id,
			destinations.name as destination_name,
			destinations.location as location,
			bookings.people_count as people_count,
			bookings.start_date as start_date,
			bookings.end_date as end_date,
			bookings.total_price as total_price,
			bookings.transport_price as transport_price,
			bookings.destination_price as destination_price,
			booking_status.name as status_name,
			payment_methods.name as payment_method_name,
			destinations.image as destination_image
		`).
		Joins("LEFT JOIN destinations ON bookings.destination_id = destinations.id").
		Joins("LEFT JOIN booking_status ON bookings.status_id = booking_status.id").
		Joins("LEFT JOIN payment_methods ON bookings.payment_method_id = payment_methods.id").
		Where("bookings.user_id = ?", userIdInt).
		Order("bookings.id DESC").
		Scan(&bookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if bookings == nil {
		bookings = []models.BookingResponse{}
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

func CancelBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	userIdInt := userID.(int)
	if booking.UserID != userIdInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only cancel your own bookings"})
		return
	}

	var status models.BookingStatus
	if err := config.DB.Table("booking_status").
		Where("LOWER(name) = ?", "canceled").
		First(&status).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Canceled status not found"})
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
		"message": "Booking canceled successfully",
		"data":    booking,
	})
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

func GetBookingById(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking

	result := config.DB.First(&booking, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Booking Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": booking,
	})
}
