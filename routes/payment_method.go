package routes

import (
	controller "github.com/Rifq11/Trava-be/controller"
	"github.com/gin-gonic/gin"
)

func PaymentMethodRoutes(router *gin.RouterGroup) {
	paymentMethods := router.Group("/payment-methods")
	{
		paymentMethods.GET("", controller.GetPaymentMethods)
	}
}

