package hello

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/lychee/foundation/logger"
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
