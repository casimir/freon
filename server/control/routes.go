package control

import (
	"net/http"
	"time"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/common"
	"github.com/casimir/freon/database"
	"github.com/casimir/freon/serialize"
	"github.com/casimir/freon/wallabag"
	"github.com/casimir/freon/wallabagproxy"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/api/tokens/schema", describer(auth.Token{}))
	r.GET("/api/tokens", func(c *gin.Context) {
		user := auth.GetUser(c)

		tokens, err := user.GetTokens()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		payload := [][]serialize.Field{}
		for _, it := range tokens {
			data, err := serialize.Describe(it)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
			payload = append(payload, data)
		}

		c.JSON(http.StatusOK, payload)
	})
	r.POST("/api/tokens", func(c *gin.Context) {
		user := auth.GetUser(c)

		var payload struct {
			Name string `json:"name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}

		if err := user.CreateToken(payload.Name); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	})
	r.GET("/api/tokens/:id", func(c *gin.Context) {
		user := auth.GetUser(c)
		id := c.Param("id")

		token, ok, err := user.GetToken(id)
		if err != nil {
			status := http.StatusNotFound
			if ok {
				status = http.StatusInternalServerError
			}
			c.AbortWithStatusJSON(status, gin.H{"message": err.Error()})
			return
		}

		payload, err := serialize.Describe(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, payload)
	})
	r.PUT("/api/tokens/:id", func(c *gin.Context) {
		user := auth.GetUser(c)
		id := c.Param("id")

		token, ok, err := user.GetToken(id)
		if err != nil {
			status := http.StatusNotFound
			if ok {
				status = http.StatusInternalServerError
			}
			c.AbortWithStatusJSON(status, gin.H{"message": err.Error()})
			return
		}

		var other auth.Token
		if err := c.ShouldBindJSON(&other); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		token.UpdateWith(&other)

		if result := database.DB.Save(&token); result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": result.Error.Error(),
			})
			return
		}

		payload, err := serialize.Describe(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, payload)
	})
	r.DELETE("/api/tokens/:id", func(c *gin.Context) {
		user := auth.GetUser(c)
		id := c.Param("id")

		ok, err := user.DeleteToken(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else if !ok {
			c.Status(http.StatusBadRequest)
		}
	})

	r.GET("/user/me", common.UserDetailHandler(nil))

	r.GET("/wallabag/credentials/schema", describer(auth.WallabagCredentials{}))
	r.GET("/wallabag/credentials", func(c *gin.Context) {
		user := auth.GetUser(c)
		if user.WallabagCredentialsID == nil {
			c.Status(http.StatusNotFound)
			return
		}

		wcreds := auth.MustGetWallabagCredentials(*user.WallabagCredentialsID)
		payload, err := serialize.Describe(wcreds)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, payload)
	})
	r.PUT("/wallabag/credentials", func(c *gin.Context) {
		user := auth.GetUser(c)

		var creds wallabag.Credentials
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if user.WallabagCredentials == nil {
			user.WallabagCredentials = &auth.WallabagCredentials{}
		}
		findErr := database.DB.Model(&user).Association("WallabagCredentials").Find(user.WallabagCredentials)
		if findErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": findErr.Error(),
			})
			return
		}
		user.WallabagCredentials.UpdateWith(&creds)

		result := database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": result.Error.Error(),
			})
			return
		}

		wcreds := user.WallabagCredentials
		client := wallabag.NewWallabagClient(wcreds.ToCredentials())
		client.FetchToken(wcreds.Username, wcreds.Password)
		wcreds.WallabagToken = client.Token()
		if err := database.DB.Save(wcreds).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	})
	r.GET("/wallabag/credentials/check", func(c *gin.Context) {
		user := auth.GetUser(c)

		wcreds := auth.WallabagCredentials{}
		err := database.DB.Model(&user).Association("WallabagCredentials").Find(&wcreds)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}

		// force the token to be refreshed so that the credentials are really checked
		wcreds.WallabagToken.ExpiresAt = time.Now().Add(-time.Second)

		_, werr := wallabagproxy.CallWallabag(&wcreds, "GET", "/api/config", nil)
		if werr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": werr.Error(),
			})
			return
		}
	})
}

func describer(v any) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := serialize.Describe(v)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, payload)
	}
}
