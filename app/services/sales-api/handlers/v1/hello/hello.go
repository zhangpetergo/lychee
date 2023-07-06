package hello

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/lychee/foundation/logger"
	log "github.com/zhangpetergo/lychee/foundation/zlog"
	"go.uber.org/zap"
	"net/http"
)

func Hello(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

func Version(c *gin.Context) {
	log := logger.NewLogger("SALES-API")

	log.Error("hello")
	c.String(http.StatusOK, "error")
}

func Test1(ctx *gin.Context) {
	name := ctx.Query("name")
	// 注意打印日志都需要通过WithContext(ctx)来获得zapLogger

	log.WithContext(ctx).Debug("测试日志", zap.String("name", name))
}
