package controller

import (
	"net/http"

	config "github.com/Rifq11/Trava-be/config"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
)

// optional i think for now
func LogActivity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.LogActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdInt := userID.(int)

	activityLog := models.ActivityLog{
		UserID:        userIdInt,
		DestinationID: req.DestinationID,
		ActivityType:  req.ActivityType,
	}

	if err := config.DB.Table("user_activity_log").Create(&activityLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Activity logged successfully",
		"data":    activityLog,
	})
}
