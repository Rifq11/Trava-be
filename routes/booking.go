package routes

import (
	controller "github.com/Rifq11/Trava-be/controller"
	middleware "github.com/Rifq11/Trava-be/middleware"
	"github.com/gin-gonic/gin"
)

func BookingRoutes(router *gin.RouterGroup) {
	bookings := router.Group("/bookings")
	bookings.Use(middleware.RequireAuth())
	{
		bookings.POST("", controller.CreateBooking)
		bookings.GET("/my", controller.GetMyBookings)
		bookings.GET("/:id", controller.GetBookingById)
		bookings.PUT("/:id/cancel", controller.CancelBooking)
		bookings.POST("/:id/cancel", controller.CancelBooking)
	}

	admin := router.Group("/admin/bookings")
	admin.Use(middleware.RequireAdmin())
	{
		admin.GET("", controller.GetAllBookingsAdmin)
		admin.GET("/:id/detail", controller.AdminGetBookingDetail)
		admin.PUT("/:id/approve", controller.ApproveBooking)
		admin.PUT("/:id/reject", controller.RejectBooking)
	}
}
 