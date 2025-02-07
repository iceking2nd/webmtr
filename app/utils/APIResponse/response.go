package APIResponse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseOKWithData(data interface{}, c *gin.Context) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    API_RESPONSE_OK,
			"data":    data,
		})
	default:
		c.String(http.StatusOK, data.(string))
	}
}

func ResponseOKWithNoContent(c *gin.Context) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusNoContent, gin.H{
			"message": "ok",
			"code":    API_RESPONSE_OK,
		})
	}
}

func ResponseError(err error, httpErrorCode, ApiErrorCode int, c *gin.Context) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(httpErrorCode, gin.H{
			"message": err.Error(),
			"code":    ApiErrorCode,
		})
	default:
		c.String(httpErrorCode, fmt.Sprintf("ERROR(%d): %v\n", ApiErrorCode, err))
	}
}
