package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhangpetergo/lychee/foundation/logger"
	"go.uber.org/zap"
	"time"
)

// 构建

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceId := uuid.NewString()
		// 加上 TraceID
		start := time.Now()

		// 原理就是
		// 1 将值 加到 zap 上
		// 2 将 加完后的 zap 加到 ctx 上
		// 3 ctx 一直传递
		logger.NewContext(c, zap.String("traceId", traceId))
		logger.NewContext(c, zap.String("method", c.Request.Method))
		logger.NewContext(c, zap.String("path", c.Request.URL.Path))
		logger.NewContext(c, zap.String("remoteaddr", c.Request.RemoteAddr))

		logger.WithContext(c).Sugar().Infow("request started", "method", c.Request.Method, "path", c.Request.URL.Path,
			"remoteaddr", c.Request.RemoteAddr)
		c.Next()
		logger.WithContext(c).Sugar().Infow("request completed", "method", c.Request.Method, "path", c.Request.URL.Path,
			"remoteaddr", c.Request.RemoteAddr, "statuscode", c.Writer.Status(), "since", time.Since(start))

	}
}
