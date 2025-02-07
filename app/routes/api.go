package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/webmtr/app/controllers/InfoController"
	"github.com/iceking2nd/webmtr/app/controllers/mtrController"
)

func apiRoutesRegister(route *gin.RouterGroup) {
	apiRoutes := route.Group("/api")

	mtrRoutes := apiRoutes.Group("/mtr/:dest_addr")
	mtrRoutes.GET("", mtrController.MTR)

	InfoRoutes := apiRoutes.Group("/info")
	InfoRoutes.GET("", InfoController.Info)
}
