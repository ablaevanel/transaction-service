package main

import (
	"log"
	"transaction-service/routes"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default()) 

	router.LoadHTMLGlob("frontend/templates/*")

	router.Static("/static", "frontend/static")

	routes.SetupRoutes(router)

	if err := router.Run(":8081"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

