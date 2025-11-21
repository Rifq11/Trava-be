package routes

import (
	controller "github.com/Rifq11/Trava-be/controller"
	helper "github.com/Rifq11/Trava-be/helper"
	middleware "github.com/Rifq11/Trava-be/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.PUT("/profile", middleware.RequireAuth(), helper.UploadSingle("user_photo"), controller.UpdateProfile)
	}
}
