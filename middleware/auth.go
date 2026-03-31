package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Store holds the cookie-backed session configuration used by the application.
var Store = sessions.NewCookieStore([]byte("wildlife-secret-key"))

// AuthRequired redirects unauthenticated requests to the login page before protected handlers run.
func AuthRequired(c *gin.Context) {
	session, err := Store.Get(c.Request, "session")
	if err != nil || session.Values["userID"] == nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Set("userID", session.Values["userID"])
	c.Next()
}
