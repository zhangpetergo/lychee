package mid

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/zhangpetergo/lychee/foundation/zlog"
	"go.uber.org/zap"
)

func TraceLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 每个请求生成的请求traceId具有全局唯一性
		traceId := uuid.NewString()
		log.NewContext(ctx, zap.String("traceId", traceId))

		// 为日志添加请求的地址以及请求参数等信息
		log.NewContext(ctx, zap.String("request.method", ctx.Request.Method))
		headers, _ := json.Marshal(ctx.Request.Header)
		log.NewContext(ctx, zap.String("request.headers", string(headers)))
		log.NewContext(ctx, zap.String("request.url", ctx.Request.URL.String()))

		// 将请求参数json序列化后添加进日志上下文
		if ctx.Request.Form == nil {
			ctx.Request.ParseMultipartForm(32 << 20)
		}
		form, _ := json.Marshal(ctx.Request.Form)
		log.NewContext(ctx, zap.String("request.params", string(form)))

		ctx.Next()
	}
}
