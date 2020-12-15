package initial

import (
	"github.com/qiaogw/pkg/redis"
	"github.com/qiaogw/pkg/store"
	"os"
	"os/signal"
	"syscall"

	"github.com/qiaogw/pkg/logs"
	_ "github.com/qiaogw/pkg/store/driver/fuse"
	_ "github.com/qiaogw/pkg/store/driver/s3"
	_ "github.com/qiaogw/pkg/store/driver/seaweed"

	"github.com/qiaogw/pkg/cache"
	"github.com/qiaogw/pkg/config"
	"github.com/qiaogw/pkg/dbtools"
	"github.com/qiaogw/pkg/jwt"
	//"github.com/qiaogw/pkg/store"
)

func init() {
	config.LoadConfig()
	if dbtools.IsInstalled(){
		if err:=dbtools.DBConnect(config.Config.DB);err!=nil{
			logs.Fatalf("jwt init err is %v",err)
		}
	}
	redis.Init(config.Config.Redis)
	cache.InitCache()
	err:=jwt.InitJwt(config.Config.JwtPublicPath,config.Config.JwtPrivatePath)
	if err!=nil{
		logs.Fatalf("jwt init err is %v",err)
	}
	caddyConfig := config.Config.Caddy
	caddyConfig.Init(config.Config.AppName)
	go caddyConfig.Start()
	store.Init()
	// 捕捉信号,配置热更新等
	go catchSignal()
	// initTask()
}

// 捕捉信号
func catchSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-c
		logs.Info("收到信号 -- ", s)
		switch s {
		case syscall.SIGHUP:
			logs.Info("收到终端断开信号, 忽略")
		case syscall.SIGINT, syscall.SIGTERM:
			shutdown()
		}
	}
}

// 应用退出
func shutdown() {
	defer func() {
		logs.Info("已退出")
		os.Exit(0)
	}()

	if !dbtools.IsInstalled() {
		return
	}
	logs.Info("应用准备退出")
	// 停止所有任务调度
	logs.Info("停止定时任务调度")
	// service.ServiceTask.WaitAndExit()
}
