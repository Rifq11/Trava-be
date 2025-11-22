package controller

import (
	"net/http"
	"strings"

	config "github.com/Rifq11/Trava-be/config"
	helper "github.com/Rifq11/Trava-be/helper"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/Rifq11/Trava-be/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleID := int(2)
	if req.RoleID != nil {
		roleID = *req.RoleID
	}

	var existingUser models.User
	result := config.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   roleID,
	}

	result = config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"data": models.RegisterResponse{
			UserID:   user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			RoleID:   user.RoleID,
			Password: user.Password,
			Token:    token,
		},
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	loginResponse := models.LoginResponse{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		RoleID:   user.RoleID,
		Password: user.Password,
		Token:    token,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    loginResponse,
	})
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIdInt := userID.(int)

	var fullName, email, phone, address, birthDate, password, userPhoto string
	contentType := c.GetHeader("Content-Type")

	if uploadedFile, ok := c.Get("uploaded_file"); ok {
		if filename, ok2 := uploadedFile.(string); ok2 && filename != "" {
			userPhoto = helper.GetFileUrl(filename)
		}
	}

	isJSONRequest := strings.Contains(contentType, "application/json")

	if userPhoto == "" && isJSONRequest {
		var req models.UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.FullName != nil {
			fullName = *req.FullName
		}
		if req.Email != nil {
			email = *req.Email
		}
		if req.Phone != nil {
			phone = *req.Phone
		}
		if req.Address != nil {
			address = *req.Address
		}
		if req.BirthDate != nil {
			birthDate = *req.BirthDate
		}
		if req.Password != nil {
			password = *req.Password
		}
	} else if userPhoto == "" {
		fullName = c.PostForm("full_name")
		email = c.PostForm("email")
		phone = c.PostForm("phone")
		address = c.PostForm("address")
		birthDate = c.PostForm("birth_date")
		password = c.PostForm("password")
	} else {
		fullName = c.PostForm("full_name")
		email = c.PostForm("email")
		phone = c.PostForm("phone")
		address = c.PostForm("address")
		birthDate = c.PostForm("birth_date")
		password = c.PostForm("password")
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

	if email != "" && email != user.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		user.Email = email
	}

	if fullName != "" {
		user.FullName = fullName
	}
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userProfile models.UserProfile
	err := config.DB.Where("user_id = ?", userIdInt).First(&userProfile).Error
	isNewProfile := false
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userProfile = models.UserProfile{UserID: userIdInt}
			isNewProfile = true
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if phone != "" {
		userProfile.Phone = phone
	}
	if address != "" {
		userProfile.Address = address
	}
	if birthDate != "" {
		userProfile.BirthDate = birthDate
	}
	if userPhoto != "" {
		userProfile.UserPhoto = userPhoto
	}

	if isNewProfile {
		if err := config.DB.Create(&userProfile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := config.DB.Save(&userProfile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	response := models.ProfileDetailResponse{
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     userProfile.Phone,
		Address:   userProfile.Address,
		BirthDate: userProfile.BirthDate,
		UserPhoto: userProfile.UserPhoto,
		Password:  user.Password,
		RoleID:    user.RoleID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    response,
	})
}
