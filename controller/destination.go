package controller

import (
	"net/http"

	config "github.com/Rifq11/Trava-be/config"
	helper "github.com/Rifq11/Trava-be/helper"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDestinations(c *gin.Context) {
	var destinations []models.DestinationResponse

	query := config.DB.
		Table("destinations").
		Select("destinations.id, destinations.name, destinations.description, destinations.location, destinations.price_per_person, destinations.image, destinations.category_id, destination_categories.name AS category_name, destinations.created_by")

	if err := query.Joins("LEFT JOIN destination_categories ON destinations.category_id = destination_categories.id").Scan(&destinations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": destinations,
	})
}

func GetDestinationById(c *gin.Context) {
	id := c.Param("id")
	var destination models.DestinationDetailResponse

	query := config.DB.
		Table("destinations").
		Select("destinations.id, destinations.name, destinations.description, destinations.location, destinations.price_per_person, destinations.image, destinations.category_id, destination_categories.name AS category_name, users.full_name AS creator_name, destinations.created_by").
		Joins("LEFT JOIN destination_categories ON destinations.category_id = destination_categories.id").
		Joins("LEFT JOIN users ON destinations.created_by = users.id").
		Where("destinations.id = ?", id)

	if err := query.First(&destination).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destination Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": destination,
	})
}

func CreateDestination(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	userIdInt := userID.(int)

	var req models.CreateDestinationRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Validation error: " + err.Error(),
		})
		return
	}

	var image string
	if uploadedFile, ok := c.Get("uploaded_file"); ok {
		if filename, ok2 := uploadedFile.(string); ok2 {
			image = helper.GetFileUrl(filename)
		}
	}
	if image == "" {
		image = req.Image
	}

	destination := models.Destination{
		CategoryID:     req.CategoryID,
		CreatedBy:      userIdInt,
		Name:           req.Name,
		Description:    req.Description,
		Location:       req.Location,
		PricePerPerson: req.PricePerPerson,
		Image:          image,
	}

	if err := config.DB.Create(&destination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Destination created successfully",
		Data:    destination,
	})
}

func UpdateDestination(c *gin.Context) {
	id := c.Param("id")

	var destination models.Destination
	if err := config.DB.First(&destination, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req models.UpdateDestinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if uploaded, exists := c.Get("uploaded_file"); exists {
		if filename, ok := uploaded.(string); ok {
			url := helper.GetFileUrl(filename)
			req.Image = &url
		}
	}

	updates := map[string]interface{}{}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.PricePerPerson != nil {
		updates["price_per_person"] = *req.PricePerPerson
	}
	if req.Image != nil {
		updates["image"] = *req.Image
	}

	if err := config.DB.Model(&destination).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Destination updated successfully",
		Data:    destination,
	})
}

func DeleteDestination(c *gin.Context) {
	id := c.Param("id")
	var destination models.Destination

	result := config.DB.First(&destination, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destination Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	result = config.DB.Delete(&destination)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Destination deleted successfully",
	})
}
