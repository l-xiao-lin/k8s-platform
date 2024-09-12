package main

import (
	"context"
	"fmt"
	"k8s-platform/controller"
	"k8s-platform/dao/mysql"
	"k8s-platform/logger"
	"k8s-platform/pkg/snowflake"
	"k8s-platform/router"
	"k8s-platform/service"
	"k8s-platform/setting"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//加载配置文件
	if err := setting.Init(); err != nil {
		fmt.Printf("setting Init failed,err:%v\n", err)
		return
	}

	//初始化zap日志
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("logger Init failed,err:%v\n", err)
		return
	}

	//初始化mysql
	if err := mysql.Init(setting.Conf.MysqlConfig); err != nil {
		fmt.Printf("mysql Init failed,err:%v\n", err)
		return
	}
	defer mysql.Close()

	//加载翻译
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("InitTrans failed,err%v\n", err)
		return
	}

	//初始化雪花算法
	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		fmt.Printf("snowflake.Init failed,err:%v\n", err)
		return
	}

	//初始化k8s clientSet
	if err := service.K8s.Init(); err != nil {
		fmt.Printf("init k8s config failed,err:%v\n", err)
		return
	}

	//优雅关机
	r := router.SetupRouter()

	srv := http.Server{Addr: fmt.Sprintf(":%d", setting.Conf.Port), Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("ListenAndServe failed,err:%v\n", err)
			return
		}
	}()

	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(fmt.Sprintf(":%d", setting.Conf.WsPort), nil)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown,err:%v\n", err)
		return
	}
	fmt.Println("Server exiting")
}
