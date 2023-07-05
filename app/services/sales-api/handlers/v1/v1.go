package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/lychee/app/services/sales-api/handlers/v1/hello"
)

// 构建 router , 处理 handle

// Routes binds all the version 1 routes.
func Routes(router *gin.Engine) {
	const version = "v1"

	//
	v1 := router.Group("/v1")
	v1.GET("/hello", hello.Hello)

}
