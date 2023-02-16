package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/routes"
	"bluebell/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err: %v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	//3.初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init myql failed, err: %v\n", err)
		return
	}
	defer mysql.Close()
	//4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err: %v\n", err)
		return
	}
	defer redis.Close()
	//5.注册路由
	r := routes.Setup()
	//6.启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.AppConfig.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
