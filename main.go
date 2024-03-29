package main

import (
	"go-rest-api/config"
	"go-rest-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
