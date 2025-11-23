package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(app *gin.Engine) {
	api := app.Group("/api")
	{
		AuthRoutes(api)
		ProfileRoutes(api)
		UserRoutes(api)
		DestinationRoutes(api)
		TransportationRoutes(api)
		BookingRoutes(api)
		PaymentRoutes(api)
		PaymentMethodRoutes(api)
		ReviewRoutes(api)
		ActivityRoutes(api)
		DashboardRoutes(api)
		ReportRoutes(api)
	}
}
