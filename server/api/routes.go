package api

import (
	"net/http"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/wallabag"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.Match([]string{http.MethodGet, http.MethodPost}, "/save", auth.WallabagAuth(), func(c *gin.Context) {
		wcreds := auth.GetWallabagCredentials(c)
		client := wallabag.NewWallabagClient(wcreds.ToCredentials())

		var payload struct {
			URL     string   `form:"url" json:"url" binding:"required"`
			Tags    []string `form:"tags[]" json:"tags"`
			Archive *bool    `form:"archive" json:"archive"`
			Starred *bool    `form:"starred" json:"starred"`
		}
		if err := c.ShouldBind(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		options := wallabag.EntriesPostOptions{
			Tags:    payload.Tags,
			Archive: payload.Archive,
			Starred: payload.Starred,
		}
		entry, err := client.Entries.Post(payload.URL, &options)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, entry)
	})

	registerUsersRoutes(r.Group("/users"))
}
