package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/casimir/freon/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	SessionDuration = 12 * time.Hour
)

const (
	CtxKeyUser                = "auth.user"
	CtxKeySession             = "auth.session"
	CtxKeyToken               = "auth.token"
	CtxKeyWallabagCredentials = "auth.wallabag_credentials"

	CookieSession = "session"
)

func findToken(value string) (*Token, error) {
	if value == "" {
		return nil, &InvalidTokenError{value, nil}
	}
	ID, err := uuid.Parse(value)
	if err != nil {
		return nil, &InvalidTokenError{value, err}
	}

	var token Token
	result := database.DB.Preload("User").Where("ID = ?", ID).Take(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, &UnknownTokenError{value}
	}
	return &token, nil
}

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("Authorization")
		if value == "" {
			value = c.Query("authorization_token")
		}

		token, err := findToken(value)
		if err != nil {
			log.Printf("token authentication failed: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(CtxKeyToken, token)
		c.Set(CtxKeyUser, &token.User)
	}
}

func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionKey, err := c.Cookie(CookieSession)
		if err != nil {
			if gin.IsDebugging() {
				log.Printf("could not get session cookie: %v", err)
			}
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		session, err := ParseSession(sessionKey)
		if err != nil {
			if gin.IsDebugging() {
				log.Printf("could not parse session: %v", err)
			}
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if time.Since(session.IssuedAt) > SessionDuration {
			if gin.IsDebugging() {
				log.Printf("session expired: %v", err)
			}
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := FindUserByID(session.UserID)
		if err != nil {
			if gin.IsDebugging() {
				log.Printf("session authentication failed: %v", err)
			}
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(CtxKeySession, session)
		c.Set(CtxKeyUser, user)
	}
}

// HardcodedAuth is a middleware that sets a hardcoded user as the authenticated user.
//
// Mostly useful for development and testing.
func HardcodedAuth(UserID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := database.DB.Take(&user, UserID).Error; err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(CtxKeyUser, &user)
	}
}

func GetUser(c *gin.Context) *User {
	user, ok := c.Get(CtxKeyUser)
	if !ok {
		panic("user *Auth middleware missing")
	}
	return user.(*User)
}

func IsSuperUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get(CtxKeyUser)
		if !ok || !user.(*User).IsSuperuser {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}

func WallabagAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)
		if user.WallabagCredentialsID == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "no wallabag credentials configured",
			})
			return
		}
		c.Set(CtxKeyWallabagCredentials, MustGetWallabagCredentials(*user.WallabagCredentialsID))
	}
}

func GetWallabagCredentials(c *gin.Context) *WallabagCredentials {
	creds, ok := c.Get(CtxKeyWallabagCredentials)
	if !ok {
		panic("WallabagAuth middleware missing")
	}
	return creds.(*WallabagCredentials)
}
