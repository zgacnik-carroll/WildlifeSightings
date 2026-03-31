package handlers

import (
	"Go-Web/db"
	"Go-Web/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AnimalCount represents the number of sightings recorded for a given animal.
type AnimalCount struct {
	Animal string
	Count  int
}

// getSessionUser returns the authenticated user's identifier and username from the session store.
func getSessionUser(c *gin.Context) (uint, string) {
	session, _ := middleware.Store.Get(c.Request, "session")
	userID, _ := session.Values["userID"].(uint)
	username, _ := session.Values["username"].(string)
	return userID, username
}

// ListSightings renders the dashboard with aggregate statistics for the entire application.
func ListSightings(c *gin.Context) {
	userID, username := getSessionUser(c)
	var totalUsers int64
	var totalSightings int64
	var counts []AnimalCount
	var topAnimal AnimalCount

	db.DB.Model(&db.User{}).Count(&totalUsers)
	db.DB.Model(&db.Sighting{}).Count(&totalSightings)
	// Build the per-animal summary shown in the dashboard table.
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

// NewSightingForm renders the form used to submit a new wildlife sighting.
func NewSightingForm(c *gin.Context) {
	_, username := getSessionUser(c)
	c.HTML(http.StatusOK, "new.html", gin.H{
		"username": username,
	})
}

// CreateSighting stores a new sighting for the authenticated user and returns to the dashboard.
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

// SearchSightings runs a simple animal/location search and renders the matching records.
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

// ShowProfile renders the signed-in user's personal sightings history.
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

// DeleteSighting removes a sighting after verifying that it belongs to the current user.
func DeleteSighting(c *gin.Context) {
	userID, _ := getSessionUser(c)
	id := c.Param("id")
	var sighting db.Sighting
	if err := db.DB.First(&sighting, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	// Enforce ownership so users cannot delete another user's records.
	if sighting.UserID != userID {
		c.Redirect(http.StatusFound, "/")
		return
	}
	db.DB.Delete(&sighting)
	c.Redirect(http.StatusFound, "/profile")
}

// ShowEditSighting loads an existing sighting into the edit form for its owner.
func ShowEditSighting(c *gin.Context) {
	userID, username := getSessionUser(c)
	id := c.Param("id")
	var sighting db.Sighting
	if err := db.DB.First(&sighting, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	// Enforce ownership so users can only edit their own sightings.
	if sighting.UserID != userID {
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"sighting": sighting,
		"username": username,
	})
}

// EditSighting updates a sighting after confirming the current user owns the record.
func EditSighting(c *gin.Context) {
	userID, _ := getSessionUser(c)
	id := c.Param("id")
	var sighting db.Sighting
	if err := db.DB.First(&sighting, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	// Enforce ownership so users cannot overwrite another user's sightings.
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
