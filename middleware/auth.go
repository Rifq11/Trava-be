package middleware

import (
	"net/http"
	"strconv"
	"strings"

	config "github.com/Rifq11/Trava-be/config"
	models "github.com/Rifq11/Trava-be/models"
	"github.com/Rifq11/Trava-be/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequireAuth middleware checks if user is authenticated
// Expects user_id in headers (x-user-id or user-id) or query params
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID int
		var err error

		authHeader := c.GetHeader("Authorization")
		if token := extractBearerToken(authHeader); token != "" {
			claims, tokenErr := utils.ParseToken(token)
			if tokenErr != nil {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Status:  "error",
					Message: "Unauthorized: Invalid token",
				})
				c.Abort()
				return
			}
			userID = claims.UserID
			c.Set("token_claims", claims)
		} else {
			userIDStr := c.GetHeader("x-user-id")
			if userIDStr == "" {
				userIDStr = c.GetHeader("user-id")
			}
			if userIDStr == "" {
				userIDStr = c.Query("user_id")
			}
			if userIDStr == "" {
				userIDStr = c.Query("userId")
			}

			if userIDStr == "" {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Status:  "error",
					Message: "Unauthorized: User ID is required",
				})
				c.Abort()
				return
			}

			userID, err = strconv.Atoi(userIDStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Status:  "error",
					Message: "Unauthorized: Invalid user ID",
				})
				c.Abort()
				return
			}
		}

		var userData struct {
			ID       int    `gorm:"column:id"`
			Email    string `gorm:"column:email"`
			FullName string `gorm:"column:full_name"`
			RoleID   int    `gorm:"column:role_id"`
			RoleName string `gorm:"column:role_name"`
		}

		if err := config.DB.
			Table("users").
			Select("users.id, users.email, users.full_name, users.role_id, roles.name as role_name").
			Joins("INNER JOIN roles ON users.role_id = roles.id").
			Where("users.id = ?", userID).
			First(&userData).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Status:  "error",
					Message: "Unauthorized: User not found",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Status:  "error",
				Message: "Authentication failed",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userData.ID)
		c.Set("user_email", userData.Email)
		c.Set("user_full_name", userData.FullName)
		c.Set("user_role_id", userData.RoleID)
		c.Set("user_role_name", userData.RoleName)

		c.Next()
	}
}

// RequireAdmin middleware checks if user is admin
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		roleName, exists := c.Get("user_role_name")
		if !exists || roleName != "admin" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Status:  "error",
				Message: "Forbidden: Requires admin role",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
