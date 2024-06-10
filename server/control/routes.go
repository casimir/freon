package control

import (
	"net/http"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/common"
	"github.com/casimir/freon/database"
	"github.com/casimir/freon/serialize"
	"github.com/casimir/freon/wallabag"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup) {
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
	r.DELETE("/api/tokens/:id", func(c *gin.Context) {
		user := auth.GetUser(c)
		id := c.Param("id")

		ok, err := user.DeleteToken(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		if !ok {
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

		var wcreds wallabag.Credentials
		if err := c.ShouldBindJSON(&wcreds); err != nil {
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
		user.WallabagCredentials.UpdateWith(&wcreds)

		result := database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": result.Error.Error(),
			})
			return
		}
	})
	r.POST("/wallabag/credentials/authenticate", func(c *gin.Context) {
		user := auth.GetUser(c)
		if user.WallabagCredentialsID == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "no wallabag credentials configured",
			})
			return
		}
		wcreds := auth.MustGetWallabagCredentials(*user.WallabagCredentialsID)

		var wuser struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBind(&wuser); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		wcreds.Username = wuser.Username
		wcreds.Password = wuser.Password

		client := wallabag.NewWallabagClient(wcreds.ToCredentials())
		client.FetchToken(wuser.Username, wuser.Password)
		wcreds.WallabagToken = client.Token()
		if err := database.DB.Save(&wcreds).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
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
