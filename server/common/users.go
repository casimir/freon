package common

import (
	"net/http"

	"github.com/casimir/freon/auth"
	"github.com/gin-gonic/gin"
)

func UserDetailHandler(getID func() string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *auth.User
		if getID == nil {
			user = auth.GetUser(c)
		} else {
			u, err := auth.FindUserByID(getID())
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			user = u
		}

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
