package routes

import (
	"net/http"
	"auth-service/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index_auth.html", nil)
	})

	router.GET("/api/auth/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	router.GET("/api/auth/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/api/auth/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "logout.html", nil)
	})

	router.POST("/api/auth/register", controllers.RegisterUser)
	router.POST("/api/auth/login", controllers.LoginUser)
}
