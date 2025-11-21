package main

import (
	config "github.com/Rifq11/Trava-be/config"
	routes "github.com/Rifq11/Trava-be/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Connect()

	r := gin.Default()

	r.Static("/uploads", "./public/uploads")

	routes.SetupRoutes(r)

	r.Run(":8080")
}
