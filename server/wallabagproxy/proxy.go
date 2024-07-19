package wallabagproxy

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/database"
	"github.com/casimir/freon/wallabag"
	"github.com/gin-gonic/gin"
)

func CallWallabag(wcreds *auth.WallabagCredentials, method string, path string, payload any) (*http.Response, error) {
	if wcreds.WallabagToken == nil {
		return nil, errors.New("no wallabag session active")
	}

	client := wallabag.NewWallabagClient(wcreds.ToCredentials())
	URL, err := client.BuildURL(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	resp, werr := client.CallAPI(method, URL, payload)

	if wcreds.WallabagToken.AccessToken != client.Token().AccessToken {
		wcreds.WallabagToken = client.Token()
		if result := database.DB.Save(wcreds); result.Error != nil {
			return nil, fmt.Errorf("failed to save token: %w", result.Error)
		}
	}

	return resp, werr
}

func RegisterRoutes(r *gin.RouterGroup) {
	r.Any("/api/*path", auth.WallabagAuth(), func(c *gin.Context) {
		wcreds := auth.GetWallabagCredentials(c)
		path := path.Join("/api", c.Param("path"))

		if len(c.Request.URL.Query()) > 0 {
			path += "?" + c.Request.URL.RawQuery
		}

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
			resp, werr = CallWallabag(wcreds, c.Request.Method, path, payload)
		} else {
			resp, werr = CallWallabag(wcreds, c.Request.Method, path, nil)
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
