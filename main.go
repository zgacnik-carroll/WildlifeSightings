package main

import (
	"Go-Web/db"
	"Go-Web/handlers"
	"Go-Web/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Public routes
	r.GET("/register", handlers.ShowRegister)
	r.POST("/register", handlers.Register)
	r.GET("/login", handlers.ShowLogin)
	r.POST("/login", handlers.Login)
	r.POST("/logout", handlers.Logout)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired)
	{
		auth.GET("/", handlers.ListSightings)
		auth.GET("/sightings/new", handlers.NewSightingForm)
		auth.POST("/sightings", handlers.CreateSighting)
		auth.GET("/sightings/search", handlers.SearchSightings)
		auth.GET("/profile", handlers.ShowProfile)
		auth.GET("/sightings/:id/edit", handlers.ShowEditSighting)
		auth.POST("/sightings/:id/edit", handlers.EditSighting)
		auth.POST("/sightings/:id/delete", handlers.DeleteSighting)
	}

	r.Run(":8080")
}
