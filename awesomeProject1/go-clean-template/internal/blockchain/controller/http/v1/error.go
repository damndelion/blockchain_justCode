package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type response struct {
	Error error `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, code int, msg error) {

	c.AbortWithStatusJSON(code, fmt.Sprintf("%s", msg))
}
