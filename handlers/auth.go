package handlers

import (
	"Go-Web/db"
	"Go-Web/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Register(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error":    "Username and password are required.",
			"username": username,
		})
		return
	}

	var existing db.User
	if err := db.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error":    "Username already taken.",
			"username": username,
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error":    "Something went wrong. Please try again.",
			"username": username,
		})
		return
	}

	user := db.User{
		Username: username,
		Password: string(hashed),
	}
	db.DB.Create(&user)
	c.Redirect(http.StatusFound, "/login")
}

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user db.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid username or password.",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Invalid username or password.",
		})
		return
	}

	session, _ := middleware.Store.Get(c.Request, "session")
	session.Values["userID"] = user.ID
	session.Values["username"] = user.Username
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	session, _ := middleware.Store.Get(c.Request, "session")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/login")
}
