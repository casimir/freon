package wallabagproxy

import (
	"io"
	"net/http"
	"path"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/database"
	"github.com/casimir/freon/wallabag"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.Any("/api/*path", auth.WallabagAuth(), func(c *gin.Context) {
		// TODO find a way to declare has route without issue with the wildcard
		if c.Param("path") == "/info" {
			infoHandler(c)
			return
		}

		wcreds := auth.GetWallabagCredentials(c)
		client := wallabag.NewWallabagClient(wcreds.ToCredentials())

		path := path.Join("/api", c.Param("path"))
		URL, _ := client.BuildURL(path, nil)

		var resp *http.Response
		var werr error
		if c.Request.ContentLength > 0 {
			payload, readErr := io.ReadAll(c.Request.Body)
			if readErr != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": readErr.Error(),
				})
				return
			}
			resp, werr = client.CallAPI(c.Request.Method, URL, payload)
		} else {
			resp, werr = client.CallAPI(c.Request.Method, URL, nil)
		}

		if wcreds.WallabagToken.AccessToken != client.Token().AccessToken {
			wcreds.WallabagToken = client.Token()
			if result := database.DB.Save(wcreds); result.Error != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": result.Error.Error(),
				})
				return
			}
		}

		if werr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": werr.Error(),
			})
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Data(resp.StatusCode, "application/json", body)
	})
}

func infoHandler(c *gin.Context) {
	client := wallabag.NewWallabagClient(wallabag.Credentials{})
	URL, _ := client.BuildURL("/api/info", nil)
	resp, werr := client.CallAPI(c.Request.Method, URL, nil)
	if werr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": werr.Error(),
		})
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.Data(resp.StatusCode, "application/json", body)
}
