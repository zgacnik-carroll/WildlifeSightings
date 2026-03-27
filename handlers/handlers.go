package handlers

import (
	"Go-Web/db"
	"Go-Web/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnimalCount struct {
	Animal string
	Count  int
}

func getSessionUser(c *gin.Context) (uint, string) {
	session, _ := middleware.Store.Get(c.Request, "session")
	userID, _ := session.Values["userID"].(uint)
	username, _ := session.Values["username"].(string)
	return userID, username
}

func ListSightings(c *gin.Context) {
	userID, username := getSessionUser(c)
	var totalUsers int64
	var totalSightings int64
	var counts []AnimalCount
	var topAnimal AnimalCount

	db.DB.Model(&db.User{}).Count(&totalUsers)
	db.DB.Model(&db.Sighting{}).Count(&totalSightings)
	db.DB.Model(&db.Sighting{}).
		Select("animal, count(*) as count").
		Group("animal").
		Order("count desc").
		Scan(&counts)

	if len(counts) > 0 {
		topAnimal = counts[0]
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"totalUsers":     totalUsers,
		"totalSightings": totalSightings,
		"counts":         counts,
		"topAnimal":      topAnimal,
		"userID":         userID,
		"username":       username,
	})
}

func NewSightingForm(c *gin.Context) {
	_, username := getSessionUser(c)
	c.HTML(http.StatusOK, "new.html", gin.H{
		"username": username,
	})
}

func CreateSighting(c *gin.Context) {
	userID, _ := getSessionUser(c)
	sighting := db.Sighting{
		Animal:   c.PostForm("animal"),
		Location: c.PostForm("location"),
		Notes:    c.PostForm("notes"),
		UserID:   userID,
	}
	db.DB.Create(&sighting)
	c.Redirect(http.StatusFound, "/")
}

func SearchSightings(c *gin.Context) {
	userID, username := getSessionUser(c)
	query := c.Query("q")
	var sightings []db.Sighting
	db.DB.Preload("User").
		Where("animal LIKE ? OR location LIKE ?", "%"+query+"%", "%"+query+"%").
		Order("created_at desc").
		Find(&sightings)
	c.HTML(http.StatusOK, "search.html", gin.H{
		"sightings": sightings,
		"query":     query,
		"userID":    userID,
		"username":  username,
	})
}

func ShowProfile(c *gin.Context) {
	userID, username := getSessionUser(c)
	var sightings []db.Sighting
	db.DB.Preload("User").Where("user_id = ?", userID).Order("created_at desc").Find(&sightings)
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"sightings": sightings,
		"username":  username,
		"userID":    userID,
	})
}

func DeleteSighting(c *gin.Context) {
	userID, _ := getSessionUser(c)
	id := c.Param("id")
	var sighting db.Sighting
	if err := db.DB.First(&sighting, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if sighting.UserID != userID {
		c.Redirect(http.StatusFound, "/")
		return
	}
	db.DB.Delete(&sighting)
	c.Redirect(http.StatusFound, "/profile")
}

func ShowEditSighting(c *gin.Context) {
	userID, username := getSessionUser(c)
	id := c.Param("id")
	var sighting db.Sighting
	if err := db.DB.First(&sighting, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if sighting.UserID != userID {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"sighting": sighting,
		"username": username,
	})
}

func EditSighting(c *gin.Context) {
	userID, _ := getSessionUser(c)
	id := c.Param("id")
	var sighting db.Sighting
	if err := db.DB.First(&sighting, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if sighting.UserID != userID {
		c.Redirect(http.StatusFound, "/")
		return
	}
	sighting.Animal = c.PostForm("animal")
	sighting.Location = c.PostForm("location")
	sighting.Notes = c.PostForm("notes")
	db.DB.Save(&sighting)
	c.Redirect(http.StatusFound, "/profile")
}
