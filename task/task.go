package task

import (
	"github.com/jasonlvhit/gocron"
	"github.com/qiaogw/log"
	"github.com/qiaogw/pkg/s3cli"
)

//task 删除超期文件
func task() {
	log.Info("开始删除超期文件。。。")
	_, count := s3cli.RemoveObjectsLife()
	log.Infof("共删除超期文件 %v 份！！！", count)
}

//taskWithParams 恢复文件
func taskWithParams(ppath string, days int64) {
	log.Info("开始恢复文件。。。")
	_, count := s3cli.RestoreObjectsLife(ppath, days)
	log.Infof("共恢复文件 %v 份！！！", count)
}

func Run() {
	//可并发运行多个任务
	//注意 interval>1时调用sAPi
	//gocron.Every(2).Seconds().Do(task)
	//gocron.Every(2).Second().Do(taskWithParams, 1, "hi")
	//在cron所有操作最后调用 start函数，否则start之后调用的操作无效不执行
	//<-gocron.Start()
	//ever := config.Config.S3.LifeDay
	//在task执行过程中 禁止异常退出
	gocron.Every(1).Day().At("02:30").DoSafely(task)
	//ppath := ""
	//days := 7
	//gocron.Every(1).Day().At("02:30").DoSafely(taskWithParams, ppath, days)
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
