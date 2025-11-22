package routes

import (
	controller "github.com/Rifq11/Trava-be/controller"
	middleware "github.com/Rifq11/Trava-be/middleware"
	"github.com/gin-gonic/gin"
)

func ReviewRoutes(router *gin.RouterGroup) {
	reviews := router.Group("/reviews")
	{
		reviews.POST("", middleware.RequireAuth(), controller.CreateReview)
		reviews.GET("/destination/:id", controller.GetDestinationReviews)
		reviews.GET("/booking/:id", middleware.RequireAuth(), controller.GetReviewByBookingID)
	}
}

