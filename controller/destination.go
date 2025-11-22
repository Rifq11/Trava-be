package controller

import (
	"net/http"
	"strconv"

	config "github.com/Rifq11/Trava-be/config"
	helper "github.com/Rifq11/Trava-be/helper"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDestinations(c *gin.Context) {
	var destinations []models.Destination

	categoryID := c.Query("category_id")
	searchQuery := c.Query("search")

	query := config.DB

	if categoryID != "" {
		categoryIDInt, err := strconv.Atoi(categoryID)
		if err == nil {
			query = query.Where("category_id = ?", categoryIDInt)
		}
	}

	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		query = query.Where("name LIKE ? OR location LIKE ?", searchPattern, searchPattern)
	}

	result := query.Find(&destinations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": destinations,
	})
}

func GetDestinationCategories(c *gin.Context) {
	var categories []models.DestinationCategory

	result := config.DB.Find(&categories)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

func GetDestinationById(c *gin.Context) {
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

	categoryIDStr := c.PostForm("category_id")
	name := c.PostForm("name")
	description := c.PostForm("description")
	location := c.PostForm("location")
	pricePerPersonStr := c.PostForm("price_per_person")

	if categoryIDStr == "" || name == "" || location == "" || pricePerPersonStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Validation error: category_id, name, location, and price_per_person are required",
		})
		return
	}

	var categoryID int
	var pricePerPerson int
	var err error

	if categoryID, err = strconv.Atoi(categoryIDStr); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Validation error: category_id must be a number",
		})
		return
	}

	if pricePerPerson, err = strconv.Atoi(pricePerPersonStr); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Validation error: price_per_person must be a number",
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
		image = c.PostForm("image")
	}

	destination := models.Destination{
		CategoryID:     categoryID,
		CreatedBy:      userIdInt,
		Name:           name,
		Description:    description,
		Location:       location,
		PricePerPerson: pricePerPerson,
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

	updates := map[string]interface{}{}

	if categoryIDStr := c.PostForm("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			updates["category_id"] = categoryID
		}
	}
	if name := c.PostForm("name"); name != "" {
		updates["name"] = name
	}
	if description := c.PostForm("description"); description != "" {
		updates["description"] = description
	}
	if location := c.PostForm("location"); location != "" {
		updates["location"] = location
	}
	if pricePerPersonStr := c.PostForm("price_per_person"); pricePerPersonStr != "" {
		if pricePerPerson, err := strconv.Atoi(pricePerPersonStr); err == nil {
			updates["price_per_person"] = pricePerPerson
		}
	}

	if uploaded, exists := c.Get("uploaded_file"); exists {
		if filename, ok := uploaded.(string); ok {
			url := helper.GetFileUrl(filename)
			updates["image"] = url
		}
	} else if image := c.PostForm("image"); image != "" {
		updates["image"] = image
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := config.DB.Model(&destination).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config.DB.First(&destination, id)

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
