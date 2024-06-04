package auth

import (
	"log"
	"net/http"

	"github.com/casimir/freon/database"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if username == "" || password == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		user, err := findUserByCredentials(username, password)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, err := NewSessionKey(user.ID.String())
		if err != nil {
			log.Printf("could not create session: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     CookieSession,
			Value:    session,
			Path:     "/",
			MaxAge:   int(SessionDuration.Seconds()),
			HttpOnly: true,
		})

		next := c.PostForm("next")
		if next == "" {
			next = "/ui/"
		}
		c.Redirect(http.StatusFound, next)
	})
	r.GET("/logout", func(c *gin.Context) {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     CookieSession,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		c.Redirect(http.StatusFound, "/ui/")
	})
}

func findUserByCredentials(username, password string) (*User, error) {
	var user User
	result := database.DB.Where("username = ?", username).Take(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, &UnknownUserError{username}
	}

	if !user.CheckPassword(password) {
		return nil, &IncorrectPasswordError{}
	}

	return &user, nil
}
