package controller

import (
	"net/http"

	config "github.com/Rifq11/Trava-be/config"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	userIdInt := userID.(int)

	var booking models.Booking
	result := config.DB.First(&booking, req.BookingID)
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

	if booking.UserID != userIdInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only review your own bookings"})
		return
	}

	review := models.Review{
		BookingID:  req.BookingID,
		UserID:     userIdInt,
		Rating:     req.Rating,
		ReviewText: req.ReviewText,
	}

	if err := config.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review created successfully",
		"data":    review,
	})
}

func GetDestinationReviews(c *gin.Context) {
	destinationID := c.Param("id")

	var reviews []models.ReviewResponse
	result := config.DB. // join to get spesific data on other table
		Table("reviews").
		Select("reviews.id, reviews.booking_id, reviews.user_id, users.full_name as user_name, reviews.rating, reviews.review_text").
		Joins("INNER JOIN bookings ON reviews.booking_id = bookings.id").
		Joins("INNER JOIN users ON reviews.user_id = users.id").
		Where("bookings.destination_id = ?", destinationID).
		Order("reviews.id DESC").
		Scan(&reviews)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reviews})
}

func GetReviewByBookingID(c *gin.Context) {
	bookingID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdInt := userID.(int)

	var review models.Review
	result := config.DB.
		Where("booking_id = ? AND user_id = ?", bookingID, userIdInt).
		First(&review)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"data": nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": review})
}
