package api

import (
	"github.com/casimir/freon/common"
	"github.com/gin-gonic/gin"
)

func registerUsersRoutes(r *gin.RouterGroup) {
	r.GET("/me", common.CurrentUserHandler())
}
