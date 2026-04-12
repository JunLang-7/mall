package router

import (
	"bytes"
	"io"
	"time"

	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type responseWriterWrapper struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *responseWriterWrapper) WriteString(s string) (int, error) {
	return w.Writer.Write([]byte(s))
}

func AccessLogMiddleware(filter func(ctx *gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}
		begin := time.Now()
		body := GetRequestBody(ctx)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))
		fields := []zap.Field{
			zap.String("ip", ctx.RemoteIP()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("params", ctx.Request.URL.RawQuery),
			zap.Any("body", body),
			zap.String("token", ctx.GetHeader(consts.UserTokenKey)),
		}
		var responseBody bytes.Buffer
		// 创建一个 MultiWriter，将响应内容同时写入原始 ResponseWriter 和 responseBody 中
		multiWriter := io.MultiWriter(ctx.Writer, &responseBody)
		// 替换原有的 ResponseWriter，使得响应内容可以被记录到日志中
		ctx.Writer = &responseWriterWrapper{
			ResponseWriter: ctx.Writer,
			Writer:         multiWriter,
		}

		ctx.Next()
		// 截断响应内容，避免日志过大
		respBody := responseBody.String()
		if len(respBody) > 1024 {
			respBody = respBody[:1024]
		}
		fields = append(fields, zap.Int64("dur_ms", time.Since(begin).Milliseconds()))
		fields = append(fields, zap.Int("status", ctx.Writer.Status()))
		fields = append(fields, zap.String("resp", respBody))
		logger.Info("access_log", fields...)
	}
}

func GetRequestBody(ctx *gin.Context) string {
	data, _ := io.ReadAll(ctx.Request.Body)
	return string(data)
}
