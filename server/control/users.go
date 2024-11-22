package control

import (
	"net/http"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/common"
	"github.com/casimir/freon/serialize"
	"github.com/gin-gonic/gin"
)

func registerUsersRoutes(r *gin.RouterGroup) {
	r.GET("/schema", describer(auth.User{}))
	r.GET("", auth.IsSuperUser(), func(c *gin.Context) {
		users, err := auth.GetAllUsers()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		payload := [][]serialize.Field{}
		for _, it := range users {
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
	r.POST("", auth.IsSuperUser(), func(c *gin.Context) {
		var payload struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}

		user, err := auth.CreateUser(payload.Username, payload.Password, false)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": user.ID})
	})
	r.GET("/:id", func(c *gin.Context) {
		user, err := auth.FindUserByID(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		common.DescribeAndRespond(c, user)
	})
	r.PUT("/:id", func(c *gin.Context) {
		c.Status(http.StatusNotImplemented)
	})
	r.DELETE("/:id", auth.IsSuperUser(), func(c *gin.Context) {
		id := c.Param("id")

		ok, err := auth.DeleteUser(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else if !ok {
			c.Status(http.StatusBadRequest)
		}
	})
	r.GET("/me", common.CurrentUserHandler())
}
