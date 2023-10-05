package main

import (
	"github.com/gin-gonic/gin"
	"invoice-view-be/controllers"
	"invoice-view-be/middleware"
	"invoice-view-be/models"
)

func main() {
	models.ConnectDatabase()

	r := gin.New()

	globalRoute := r.Use(middleware.CORSMiddleware())

	globalRoute.POST("/api/login", controllers.Login)
	globalRoute.POST("/api/register", controllers.Register)

	protected := globalRoute.Use(middleware.JwtAuthMiddleWare())
	protected.GET("/api/user", controllers.CurrentUser)

	r.Run()
}
