package handlers

import (
	"Go-Web/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListSightings(c *gin.Context) {
	var sightings []db.Sighting
	db.DB.Order("created_at desc").Find(&sightings)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"sightings": sightings,
	})
}

func NewSightingForm(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", nil)
}

func CreateSighting(c *gin.Context) {
	sighting := db.Sighting{
		Animal:   c.PostForm("animal"),
		Location: c.PostForm("location"),
		Notes:    c.PostForm("notes"),
	}
	db.DB.Create(&sighting)
	c.Redirect(http.StatusFound, "/")
}

func SearchSightings(c *gin.Context) {
	query := c.Query("q")
	var sightings []db.Sighting
	db.DB.Where("animal LIKE ? OR location LIKE ?", "%"+query+"%", "%"+query+"%").
		Order("created_at desc").
		Find(&sightings)
	c.HTML(http.StatusOK, "search.html", gin.H{
		"sightings": sightings,
		"query":     query,
	})
}

func GetStats(c *gin.Context) {
	type AnimalCount struct {
		Animal string
		Count  int
	}

	var counts []AnimalCount
	db.DB.Model(&db.Sighting{}).
		Select("animal, count(*) as count").
		Group("animal").
		Order("count desc").
		Scan(&counts)

	var total int64
	db.DB.Model(&db.Sighting{}).Count(&total)

	c.HTML(http.StatusOK, "stats.html", gin.H{
		"counts": counts,
		"total":  total,
	})
}
