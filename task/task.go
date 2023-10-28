package task

import (
	"os/exec"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/qiaogw/pkg/logs"
)

func init() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	gocron.ChangeLoc(location)
}

func Run() {
	// 配置本地时间
	//location, _ := time.LoadLocation("Asia/Shanghai")
	//gocron.ChangeLoc(location)
	//可并发运行多个任务
	//注意 interval>1时调用sAPi
	//gocron.Every(2).Seconds().Do(task)
	//gocron.Every(2).Second().Do(taskWithParams, 1, "hi")
	//在cron所有操作最后调用 start函数，否则start之后调用的操作无效不执行
	//<-gocron.Start()
	//ever := config.Config.S3.LifeDay
	//在task执行过程中 禁止异常退出
	//gocron.Every(1).Day().At(config.Config.S3.TaskTime).DoSafely(task)
	//gocron.Every(1).Day().Do(task)
	//ppath := config.Config.S3.DefautRestorePath
	//days := 7
	//gocron.Every(1).Day().At("03:30").DoSafely(taskWithParams, ppath, days)
	//
	//// 删除某一任务
	//gocron.Remove(task)
	//
	////删除所有任务
	//gocron.Clear()

	//可同时创建一个新的任务调度 2个schedulers 同时执行
	//s := gocron.NewScheduler()
	//s.Every(3).Seconds().Do(task)
	//<-s.Start()

	//防止多个集群中任务同时执行 task 实现lock接口
	//两行代码，对cron 设置lock实现，执行task时调用Lock方法再Do task
	//gocron.SetLocker(lockerImplementation)
	//gocron.Every(1).Hour().Lock().Do(task)

	<-gocron.Start()
}

func AddTask(name string, args ...string) {
	// 配置本地时间

	gocron.Every(1).Day().At("03:30").DoSafely(runCmd, name, args)
	<-gocron.Start()
}

func runCmd(name string, args ...string) {
	command := exec.Command(name, args...)
	if err := command.Start(); err != nil {
		logs.Error("Failed to start cmd: ", err)
	}
	if err := command.Wait(); err != nil {
		logs.Error("Failed to Wait cmd: ", err) //
	}
}
