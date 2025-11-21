package controller

import (
	"net/http"
	"strings"

	config "github.com/Rifq11/Trava-be/config"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User

	result := config.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if users == nil {
		users = []models.User{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	result := config.DB.First(&user, id)
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

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
	})
}

func CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	result := config.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	roleID := int(2)
	if req.RoleID != nil {
		roleID = *req.RoleID
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

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	result := config.DB.First(&user, id)
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

	if req.Email != nil && *req.Email != user.Email {
		var existingUser models.User
		result := config.DB.Where("email = ?", *req.Email).First(&existingUser)
		if result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}
	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}

	result = config.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	result := config.DB.First(&user, id)
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

	result = config.DB.Delete(&user)
	if result.Error != nil {
		errStr := result.Error.Error()
		if strings.Contains(errStr, "1451") || strings.Contains(errStr, "foreign key constraint") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Cannot delete user: user is referenced by other records",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
