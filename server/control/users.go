package control

import (
	"net/http"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/common"
	"github.com/casimir/freon/database"
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
		id := c.Param("id")
		currentUser := auth.GetUser(c)
		if !(currentUser.ID.String() == id || currentUser.IsSuperuser) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		user, err := auth.FindUserByID(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if user == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		var payload struct {
			Username     string `json:"username"`
			OldPassword  string `json:"old_password"`
			NewPassword1 string `json:"new_password1"`
			NewPassword2 string `json:"new_password2"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		}

		dirty := false
		if payload.Username != "" {
			user.Username = payload.Username
			dirty = true
		}
		if payload.OldPassword != "" || payload.NewPassword1 != "" || payload.NewPassword2 != "" {
			if payload.NewPassword1 != payload.NewPassword2 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "new passwords do not match",
				})
				return
			}
			if !user.CheckPassword(payload.OldPassword) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "incorrect password",
				})
				return
			}
			user.SetPassword(payload.NewPassword1)
			dirty = true
		}

		if dirty {
			if err := database.DB.Save(user).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
			c.Status(http.StatusOK)
		}
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
