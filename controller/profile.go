package controller

import (
	"net/http"

	config "github.com/Rifq11/Trava-be/config"
	helper "github.com/Rifq11/Trava-be/helper"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	userIdInt := userID.(int)

	var user models.User
	result := config.DB.First(&user, userIdInt)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	var userProfile models.UserProfile
	var adminProfile models.AdminProfile
	var profile interface{}

	if result := config.DB.Where("user_id = ?", userIdInt).First(&userProfile); result.Error == nil {
		profile = userProfile
	} else if result := config.DB.Where("user_id = ?", userIdInt).First(&adminProfile); result.Error == nil {
		profile = adminProfile
	}

	response := models.ProfileResponse{
		User: models.ProfileUserResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Password: user.Password,
			RoleID:   user.RoleID,
		},
		Profile: profile,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func CompleteProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdInt := userID.(int)

	var userPhoto string
	if uploadedFile, exists := c.Get("uploaded_file"); exists {
		if filename, ok := uploadedFile.(string); ok {
			// get url
			userPhoto = helper.GetFileUrl(filename)
		}
	}
	if userPhoto == "" {
		userPhoto = c.PostForm("user_photo")
		if userPhoto == "" {
			userPhoto = c.PostForm("userPhoto")
		}
	}

	phone := c.PostForm("phone")
	address := c.PostForm("address")
	birthDate := c.PostForm("birth_date")
	if birthDate == "" {
		birthDate = c.PostForm("birthDate")
	}

	var req models.CompleteProfileRequest
	if phone != "" {
		req.Phone = &phone
	}
	if address != "" {
		req.Address = &address
	}
	if birthDate != "" {
		req.BirthDate = &birthDate
	}
	if userPhoto != "" {
		req.UserPhoto = &userPhoto
	}
	roleName, _ := c.Get("user_role_name")
	roleIDVal, _ := c.Get("user_role_id")
	var adminProfile models.AdminProfile
	var userProfile models.UserProfile

	var user models.User
	result := config.DB.First(&user, userIdInt)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	isAdmin := false
	if roleNameStr, ok := roleName.(string); ok && roleNameStr == "admin" {
		isAdmin = true
	}
	if !isAdmin {
		if roleIDInt, ok := roleIDVal.(int); ok && roleIDInt == 1 {
			isAdmin = true
		}
	}

	if isAdmin {
		result := config.DB.Where("user_id = ?", userIdInt).First(&adminProfile)
		if result.Error == nil {
			if req.Phone != nil {
				adminProfile.Phone = *req.Phone
			}
			if req.Address != nil {
				adminProfile.Address = *req.Address
			}
			if req.BirthDate != nil {
				adminProfile.BirthDate = *req.BirthDate
			}
			if req.UserPhoto != nil {
				adminProfile.UserPhoto = *req.UserPhoto
			}
			adminProfile.IsCompleted = true
			if err := config.DB.Save(&adminProfile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			newProfile := models.AdminProfile{
				UserID:      userIdInt,
				IsCompleted: true,
			}
			if req.Phone != nil {
				newProfile.Phone = *req.Phone
			}
			if req.Address != nil {
				newProfile.Address = *req.Address
			}
			if req.BirthDate != nil {
				newProfile.BirthDate = *req.BirthDate
			}
			if req.UserPhoto != nil {
				newProfile.UserPhoto = *req.UserPhoto
			}
			if err := config.DB.Create(&newProfile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			adminProfile = newProfile
		}
	} else {
		result := config.DB.Where("user_id = ?", userIdInt).First(&userProfile)
		if result.Error == nil {
			if req.Phone != nil {
				userProfile.Phone = *req.Phone
			}
			if req.Address != nil {
				userProfile.Address = *req.Address
			}
			if req.BirthDate != nil {
				userProfile.BirthDate = *req.BirthDate
			}
			if req.UserPhoto != nil {
				userProfile.UserPhoto = *req.UserPhoto
			}
			userProfile.IsCompleted = true
			if err := config.DB.Save(&userProfile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			newProfile := models.UserProfile{
				UserID:      userIdInt,
				IsCompleted: true,
			}
			if req.Phone != nil {
				newProfile.Phone = *req.Phone
			}
			if req.Address != nil {
				newProfile.Address = *req.Address
			}
			if req.BirthDate != nil {
				newProfile.BirthDate = *req.BirthDate
			}
			if req.UserPhoto != nil {
				newProfile.UserPhoto = *req.UserPhoto
			}
			if err := config.DB.Create(&newProfile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			userProfile = newProfile
		}
	}

	var response models.ProfileDetailResponse
	if isAdmin {
		response = models.ProfileDetailResponse{
			FullName:  user.FullName,
			Email:     user.Email,
			Phone:     adminProfile.Phone,
			Address:   adminProfile.Address,
			BirthDate: adminProfile.BirthDate,
			UserPhoto: adminProfile.UserPhoto,
			Password:  user.Password,
			RoleID:    user.RoleID,
		}
	} else {
		response = models.ProfileDetailResponse{
			FullName:  user.FullName,
			Email:     user.Email,
			Phone:     userProfile.Phone,
			Address:   userProfile.Address,
			BirthDate: userProfile.BirthDate,
			UserPhoto: userProfile.UserPhoto,
			Password:  user.Password,
			RoleID:    user.RoleID,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile completed successfully",
		"data":    response,
	})
}
