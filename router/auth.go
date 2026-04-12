package router

import (
	"context"
	"net/http"

	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/gin-gonic/gin"
)

type TokenFunc func(ctx context.Context, token string) (*common.User, error)
type TokenAdminFunc func(ctx context.Context, token string) (*common.AdminUser, error)

// AuthMiddleware 用户侧鉴权中间件
func AuthMiddleware(filter func(ctx *gin.Context) bool, getTokenFunc TokenFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}
		token := ctx.GetHeader(consts.UserTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		user, err := getTokenFunc(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		ctx.Set(consts.CustomerUserKey, user)
		ctx.Next()
	}
}

// AdminAuthMiddleware 管理后台用户侧鉴权中间件
func AdminAuthMiddleware(filter func(ctx *gin.Context) bool, getTokenFunc TokenAdminFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}
		token := ctx.GetHeader(consts.AdminTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		admin, err := getTokenFunc(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		ctx.Set(consts.AdminUserKey, admin)
		ctx.Next()
	}
}
