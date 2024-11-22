package common

import (
	"net/http"

	"github.com/casimir/freon/auth"
	"github.com/gin-gonic/gin"
)

func CurrentUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := auth.GetUser(c)

		payload := struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			IsSuperuser bool   `json:"is_superuser"`
		}{
			ID:          user.ID.String(),
			Username:    user.Username,
			IsSuperuser: user.IsSuperuser,
		}
		c.JSON(http.StatusOK, &payload)
	}
}
