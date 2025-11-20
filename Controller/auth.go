package controller

import (
	"net/http"

	config "github.com/Rifq11/Trava-be/Config"
	models "github.com/Rifq11/Trava-be/Models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Full name, email, and password are required",
		})
		return
	}

	roleID := int(2)
	if req.RoleID != nil {
		roleID = *req.RoleID
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Status:  "error",
			Message: "User with this email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to register user",
		})
		return
	}

	user := models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   roleID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to register user",
		})
		return
	}

	registerResponse := models.RegisterResponse{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: user.FullName,
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    registerResponse,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Email and password are required",
		})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Status:  "error",
				Message: "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to login",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid email or password",
		})
		return
	}

	var role models.Role
	if err := config.DB.First(&role, user.RoleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to login",
		})
		return
	}

	loginResponse := models.LoginResponse{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		RoleID:   user.RoleID,
		RoleName: role.Name,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    loginResponse,
	})
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Check if user exists
	var user models.User
	userIdInt := userID.(int)
	if err := config.DB.First(&user, userIdInt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Status:  "error",
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to update profile",
		})
		return
	}

	if req.Email != nil && *req.Email != user.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", *req.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Status:  "error",
				Message: "User with this email already exists",
			})
			return
		}
		user.Email = *req.Email
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Status:  "error",
				Message: "Failed to update profile",
			})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to update profile",
		})
		return
	}

	var userProfile models.UserProfile
	if err := config.DB.Where("user_id = ?", userIdInt).First(&userProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newProfile := models.UserProfile{
				UserID:      userIdInt,
				Phone:       "",
				Address:     "",
				BirthDate:   "",
				IsCompleted: false,
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
			config.DB.Create(&newProfile)
		}
	} else {
		if req.Phone != nil {
			userProfile.Phone = *req.Phone
		}
		if req.Address != nil {
			userProfile.Address = *req.Address
		}
		if req.BirthDate != nil {
			userProfile.BirthDate = *req.BirthDate
		}
		config.DB.Save(&userProfile)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  "success",
		Message: "Profile updated successfully",
	})
}
