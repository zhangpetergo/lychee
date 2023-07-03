package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		time.Sleep(10 * time.Second)
		c.String(http.StatusOK, "hello")
	})

	api := http.Server{
		Addr:    ":9999",
		Handler: router,
	}

	// 接收关闭信号
	shutdown := make(chan os.Signal, 1)

	// kill 没有参数时发送 syscall.SIGTERM
	// kill -2 ,也就是 ctrl + c ,系统发送 syscall.SIGINT
	// kill -9 发送 SIGKILL,SIGKILL是一个无法被捕获、阻塞或忽略的终止信号
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverError := make(chan error, 1)

	go func() {
		log.Println("俺启动了")
		if err := api.ListenAndServe(); err != nil {
			serverError <- err
		}
	}()

	select {
	case err := <-serverError:
		log.Println("server error: ", err)
	case sig := <-shutdown:
		log.Printf("receive signal:%v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := api.Shutdown(ctx); err != nil {
			log.Println("shutdown:", err)
			api.Close()
		}
		// 做一些资源关闭，清理操作
	}
}
