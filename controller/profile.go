package controller

import (
	"net/http"
	"strings"

	config "github.com/Rifq11/Trava-be/config"
	helper "github.com/Rifq11/Trava-be/helper"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdInt := userID.(int)

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

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func CompleteProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdInt := userID.(int)

	// upload photo (optional)
	var userPhoto string
	if uploadedFile, exists := c.Get("uploaded_file"); exists {
		if filename, ok := uploadedFile.(string); ok {
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

	if birthDate != "" {
		if strings.Contains(birthDate, "T") {
			birthDate = strings.Split(birthDate, "T")[0]
		}
		birthDate = strings.TrimSuffix(birthDate, "Z")
		if len(birthDate) > 10 {
			birthDate = birthDate[:10]
		}
	}

	updates := map[string]interface{}{
		"is_completed": true,
	}

	if phone != "" {
		updates["phone"] = phone
	}
	if address != "" {
		updates["address"] = address
	}
	if birthDate != "" {
		updates["birth_date"] = birthDate
	}
	if userPhoto != "" {
		updates["user_photo"] = userPhoto
	}

	roleName, _ := c.Get("user_role_name")
	roleIDVal, _ := c.Get("user_role_id")

	isAdmin := false
	if roleNameStr, ok := roleName.(string); ok && roleNameStr == "admin" {
		isAdmin = true
	}
	if !isAdmin {
		if roleIDInt, ok := roleIDVal.(int); ok && roleIDInt == 1 {
			isAdmin = true
		}
	}

	var adminProfile models.AdminProfile
	var userProfile models.UserProfile

	if isAdmin {
		result := config.DB.Where("user_id = ?", userIdInt).First(&adminProfile)

		if result.Error == nil {
			// update
			if err := config.DB.Model(&adminProfile).Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			// create
			newProfile := models.AdminProfile{
				UserID:      userIdInt,
				IsCompleted: true,
			}

			if phone != "" {
				newProfile.Phone = phone
			}
			if address != "" {
				newProfile.Address = address
			}
			if birthDate != "" {
				newProfile.BirthDate = birthDate
			}
			if userPhoto != "" {
				newProfile.UserPhoto = userPhoto
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
			if err := config.DB.Model(&userProfile).Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			newProfile := models.UserProfile{
				UserID:      userIdInt,
				IsCompleted: true,
			}

			if phone != "" {
				newProfile.Phone = phone
			}
			if address != "" {
				newProfile.Address = address
			}
			if birthDate != "" {
				newProfile.BirthDate = birthDate
			}
			if userPhoto != "" {
				newProfile.UserPhoto = userPhoto
			}

			if err := config.DB.Create(&newProfile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			userProfile = newProfile
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile completed successfully"})
}
