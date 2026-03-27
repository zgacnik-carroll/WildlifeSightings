package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("wildlife-secret-key"))

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
