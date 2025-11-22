package routes

import (
	controller "github.com/Rifq11/Trava-be/controller"
	"github.com/gin-gonic/gin"
)

func TransportationRoutes(router *gin.RouterGroup) {
	transportations := router.Group("/transportations")
	{
		transportations.GET("/destination/:id", controller.GetTransportationsByDestination)
	}
}

