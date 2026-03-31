package handlers

import (
	"Go-Web/db"
	"Go-Web/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ShowRegister renders the registration form.
func ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// Register validates the submitted credentials, creates the user record, and redirects to login.
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

	// Hash the password before persisting it so plaintext credentials are never stored.
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

// ShowLogin renders the login form.
func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Login authenticates the user and stores the session identifiers in the cookie-backed session.
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
	// Persist the user identity in the session so protected routes can authorize the request.
	session.Values["userID"] = user.ID
	session.Values["username"] = user.Username
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/")
}

// Logout clears the current session and sends the user back to the login screen.
func Logout(c *gin.Context) {
	session, _ := middleware.Store.Get(c.Request, "session")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/login")
}
