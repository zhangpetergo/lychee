package hello

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/lychee/foundation/logger"
	"go.uber.org/zap"
	"net/http"
)

func Hello(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

func Version(c *gin.Context) {
	logger.WithContext(c).Error("error")
	c.String(http.StatusOK, "error")
}

func Test1(ctx *gin.Context) {
	name := ctx.Query("name")
	// 注意打印日志都需要通过WithContext(ctx)来获得zapLogger

	logger.WithContext(ctx).Debug("测试日志", zap.String("name", name))
}
