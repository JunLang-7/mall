package router

import "github.com/gin-gonic/gin"

func AccessLogMiddleware(filter func(ctx *gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}
		ctx.Next()
	}
}
