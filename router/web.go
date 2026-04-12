package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/JunLang-7/mall/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	server *gin.Engine
	addr   string
}

func NewApp(port int, router IRouter) *App {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	// Recover middleware: 全局捕获panic
	engine.Use(gin.Recovery())
	// 日志中间件，自定义过滤器
	engine.Use(AccessLogMiddleware(router.AccessRecordFilter))
	// 注册业务路由
	router.Register(engine)
	return &App{
		server: engine,
		addr:   ":" + strconv.Itoa(port),
	}
}

func (a *App) Run() {
	srv := &http.Server{
		Addr:    a.addr,
		Handler: a.server,
	}

	// 异步启动
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen err: %v", err)
		}
	}()

	logger.Debug(fmt.Sprintf("start server, endpoint: http://localhost%s", a.addr))
	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	msg := <-closeCh
	logger.Warn("server shutdown", zap.String("msg", msg.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown err: %v", err)
	}
}
