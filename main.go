package main

import (
	"Go-Web/db"
	"Go-Web/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", handlers.ListSightings)
	r.GET("/sightings/new", handlers.NewSightingForm)
	r.POST("/sightings", handlers.CreateSighting)
	r.GET("/sightings/search", handlers.SearchSightings)
	r.GET("/stats", handlers.GetStats)

	r.Run(":8080")
}
