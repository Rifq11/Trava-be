package routes

import (
	controller "github.com/Rifq11/Trava-be/controller"
	middleware "github.com/Rifq11/Trava-be/middleware"
	"github.com/gin-gonic/gin"
)

func ReportRoutes(router *gin.RouterGroup) {
	reports := router.Group("/reports")
	reports.Use(middleware.RequireAdmin())
	reports.GET("/orders", controller.GetReportOrders)
	reports.GET("/income", controller.GetIncomeReport)
}
