package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/webmtr/global"
)

func SetupRouter(router *gin.RouterGroup) {
	apiRoutesRegister(router)
	if global.LogLevel >= 5 {
		debugRoutesRegister(router)
	}
}
