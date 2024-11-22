package common

import (
	"net/http"

	"github.com/casimir/freon/serialize"
	"github.com/gin-gonic/gin"
)

func DescribeAndRespond(c *gin.Context, obj any) {
	payload, err := serialize.Describe(obj)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, payload)
}
