package main

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/zhangpetergo/lychee/app/services/sales-api/handlers"
	"github.com/zhangpetergo/lychee/business/web/v1/debug"
	"github.com/zhangpetergo/lychee/foundation/logger"
	_ "github.com/zhangpetergo/lychee/foundation/zlog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// =========================================================================
	// 初始化日志

	log := logger.NewLogger("SALES-API").Sugar()
	defer log.Sync()

	// =========================================================================
	// 读取配置文件

	// 通过读取环境变量切换不同的配置文件
	viper.AutomaticEnv()
	test := viper.GetBool("TEST")

	// 默认为dev
	configFileName := "config.dev.yaml"
	if test {
		configFileName = "config.test.yaml"
	}
	config := new(Config)
	viper.SetConfigFile(configFileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		panic(fmt.Sprintf("read config fail %v \n", err))
	}

	// 将配置文件映射到结构体
	if err = viper.Unmarshal(config); err != nil {
		log.Error(err)
		panic(fmt.Sprintf("read config fail %v \n", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()
	// 注意！！！配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Infow("修改配置文件")
		if err := viper.Unmarshal(config); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})

	// =========================================================================
	// Start Debug Service

	go func() {
		log.Infow("Service Debug Start", "DebugHost", config.Server.DebugHost)
		if err := http.ListenAndServe(config.Server.DebugHost, debug.StandardLibraryMux()); err != nil {
			log.Errorw("Service Debug Shutdown", "HOST", config.Server.DebugHost)
		}
	}()

	api := http.Server{
		Addr:    config.Server.APIHost,
		Handler: handlers.Router(log),
	}

	// 接收关闭信号
	shutdown := make(chan os.Signal, 1)

	// kill 没有参数时发送 syscall.SIGTERM
	// kill -2 ,也就是 ctrl + c ,系统发送 syscall.SIGINT
	// kill -9 发送 SIGKILL,SIGKILL是一个无法被捕获、阻塞或忽略的终止信号
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverError := make(chan error, 1)

	go func() {
		log.Infow("Service Start", "ApiHost", config.Server.APIHost)
		if err := api.ListenAndServe(); err != nil {
			serverError <- err
		}
	}()

	select {
	case err := <-serverError:
		log.Info("server error: ", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := api.Shutdown(ctx); err != nil {
			log.Info("shutdown:", err)
			api.Close()
		}
		// 做一些资源关闭，清理操作
	}
}
