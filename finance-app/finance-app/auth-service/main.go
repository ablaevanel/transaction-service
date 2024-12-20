package main

import (
	"log"
	"auth-service/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.LoadHTMLGlob("/app/frontend/templates/*")

	router.Static("/static", "/app/frontend/static")

	routes.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
