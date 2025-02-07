package routes

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

func debugRoutesRegister(route *gin.RouterGroup) {
	debugRoutes := route.Group("/debug")
	debugRoutes.GET("/pprof/", func(ctx *gin.Context) { pprof.Index(ctx.Writer, ctx.Request) })
	debugRoutes.GET("/pprof/:1", func(ctx *gin.Context) { pprof.Index(ctx.Writer, ctx.Request) })
	debugRoutes.GET("/pprof/trace", func(ctx *gin.Context) { pprof.Trace(ctx.Writer, ctx.Request) })
	debugRoutes.GET("/pprof/symbol", func(ctx *gin.Context) { pprof.Symbol(ctx.Writer, ctx.Request) })
	debugRoutes.GET("/pprof/profile", func(ctx *gin.Context) { pprof.Profile(ctx.Writer, ctx.Request) })
	debugRoutes.GET("/pprof/cmdline", func(ctx *gin.Context) { pprof.Cmdline(ctx.Writer, ctx.Request) })
}
