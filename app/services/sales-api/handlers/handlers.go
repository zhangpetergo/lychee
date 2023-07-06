package handlers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	v1 "github.com/zhangpetergo/lychee/app/services/sales-api/handlers/v1"
	"github.com/zhangpetergo/lychee/business/web/v1/mid"
	"go.uber.org/zap"
	"net/http"
)

// 这个返回一个 router

func Router(log *zap.SugaredLogger) http.Handler {
	// 初始化 router
	router := gin.New()
	// 中间件
	// 捕获 panic
	router.Use(gin.Recovery())
	// 压缩
	router.Use(gzip.Gzip(gzip.BestCompression))
	// 日志跟踪
	//router.Use(mid.Logger(log))
	router.Use(mid.Logger())
	v1.Routes(router)
	return router
}
