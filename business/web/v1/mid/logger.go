package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// 构建

func Logger(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {

		traceID := uuid.NewString()
		// 加上 TraceID
		start := time.Now()
		log.Infow("request started", "traceID", traceID, "method", c.Request.Method, "path", c.Request.URL.Path,
			"remoteaddr", c.Request.RemoteAddr)
		c.Next()
		log.Infow("request completed", "traceID", traceID, "method", c.Request.Method, "path", c.Request.URL.Path,
			"remoteaddr", c.Request.RemoteAddr, "statuscode", c.Writer.Status(), "since", time.Since(start))

	}
}
