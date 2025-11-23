package routes

import (
	controller "github.com/Rifq11/Trava-be/controller"
	middleware "github.com/Rifq11/Trava-be/middleware"
	"github.com/gin-gonic/gin"
)

func DashboardRoutes(router *gin.RouterGroup) {
	dashboard := router.Group("/dashboard")
	dashboard.Use(middleware.RequireAdmin())
	{
		dashboard.GET("", controller.GetDashboardStatistics)
		dashboard.GET("/monthly-sales", controller.GetMonthlySales)
	}
}