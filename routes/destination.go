package routes

import (
	controller "github.com/Rifq11/Trava-be/controller"
	helper "github.com/Rifq11/Trava-be/helper"
	middleware "github.com/Rifq11/Trava-be/middleware"
	"github.com/gin-gonic/gin"
)

func DestinationRoutes(router *gin.RouterGroup) {
	destinations := router.Group("/destinations")
	{
		destinations.GET("", controller.GetDestinations)
		destinations.GET("/categories", controller.GetDestinationCategories)
		destinations.GET("/:id", controller.GetDestinationById)
		destinations.POST("", middleware.RequireAuth(), helper.UploadSingle("image"), controller.CreateDestination)
		destinations.PUT("/:id", helper.UploadSingle("image"), controller.UpdateDestination)
		destinations.DELETE("/:id", controller.DeleteDestination)
	}
}
