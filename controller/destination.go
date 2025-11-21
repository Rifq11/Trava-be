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
	categoryIDStr := c.Query("category_id")
	var destinations []models.DestinationResponse

	query := config.DB.
		Table("destinations").
		Select("destinations.id, destinations.name, destinations.description, destinations.location, destinations.price_per_person, destinations.image, destinations.category_id, destinations.created_by, destination_categories.name as category_name").
		Joins("INNER JOIN destination_categories ON destinations.category_id = destination_categories.id").
		Order("destinations.id DESC")

	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			query = query.Where("destinations.category_id = ?", categoryID)
		}
	}

	if err := query.Scan(&destinations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to get destinations",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status: "success",
		Data:   destinations,
	})
}

func GetDestinationById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid destination ID",
		})
		return
	}

	var destination models.DestinationDetailResponse
	if err := config.DB.
		Table("destinations").
		Select("destinations.id, destinations.name, destinations.description, destinations.location, destinations.price_per_person, destinations.image, destinations.category_id, destinations.created_by, destination_categories.name as category_name, users.full_name as creator_name").
		Joins("INNER JOIN destination_categories ON destinations.category_id = destination_categories.id").
		Joins("INNER JOIN users ON destinations.created_by = users.id").
		Where("destinations.id = ?", id).
		First(&destination).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Status:  "error",
				Message: "Destination not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to get destination",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status: "success",
		Data:   destination,
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

	var image string
	if uploadedFile, exists := c.Get("uploaded_file"); exists {
		if filename, ok := uploadedFile.(string); ok {
			// get url
			image = helper.GetFileUrl(filename)
		}
	}
	if image == "" {
		image = c.PostForm("image")
	}

	categoryIDStr := c.PostForm("category_id")
	if categoryIDStr == "" {
		categoryIDStr = c.PostForm("categoryId")
	}
	name := c.PostForm("name")
	description := c.PostForm("description")
	location := c.PostForm("location")
	pricePerPersonStr := c.PostForm("price_per_person")
	if pricePerPersonStr == "" {
		pricePerPersonStr = c.PostForm("pricePerPerson")
	}

	if categoryIDStr == "" || name == "" || location == "" || pricePerPersonStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Category ID, name, location, and price per person are required",
		})
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid category ID",
		})
		return
	}

	pricePerPerson, err := strconv.Atoi(pricePerPersonStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid price per person",
		})
		return
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
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to create destination",
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Status:  "success",
		Message: "Destination created successfully",
		Data:    map[string]interface{}{"id": destination.ID},
	})
}

func UpdateDestination(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid destination ID",
		})
		return
	}

	var destination models.Destination
	if err := config.DB.First(&destination, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Status:  "error",
				Message: "Destination not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to update destination",
		})
		return
	}

	var image string
	if uploadedFile, exists := c.Get("uploaded_file"); exists {
		if filename, ok := uploadedFile.(string); ok {
			// get url
			image = helper.GetFileUrl(filename)
		}
	}
	if image == "" {
		image = c.PostForm("image")
	}

	categoryIDStr := c.PostForm("category_id")
	if categoryIDStr == "" {
		categoryIDStr = c.PostForm("categoryId")
	}
	name := c.PostForm("name")
	description := c.PostForm("description")
	location := c.PostForm("location")
	pricePerPersonStr := c.PostForm("price_per_person")
	if pricePerPersonStr == "" {
		pricePerPersonStr = c.PostForm("pricePerPerson")
	}

	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			destination.CategoryID = categoryID
		}
	}
	if name != "" {
		destination.Name = name
	}
	if description != "" {
		destination.Description = description
	}
	if location != "" {
		destination.Location = location
	}
	if pricePerPersonStr != "" {
		pricePerPerson, err := strconv.Atoi(pricePerPersonStr)
		if err == nil {
			destination.PricePerPerson = pricePerPerson
		}
	}
	if image != "" {
		destination.Image = image
	}

	if err := config.DB.Save(&destination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to update destination",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  "success",
		Message: "Destination updated successfully",
	})
}

func DeleteDestination(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid destination ID",
		})
		return
	}

	var destination models.Destination
	if err := config.DB.First(&destination, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Status:  "error",
				Message: "Destination not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete destination",
		})
		return
	}

	if err := config.DB.Delete(&destination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete destination",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  "success",
		Message: "Destination deleted successfully",
	})
}
