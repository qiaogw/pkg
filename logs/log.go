// Copyright 2018 cloudy itcloudy@qq.com.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

type LogKafka struct {
	Producer sarama.SyncProducer
	Topic    string
}

func (lk *LogKafka) Write(p []byte) (n int, err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = lk.Topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err = lk.Producer.SendMessage(msg)
	if err != nil {
		return
	}
	return
}
func InitLogger(mode string, fileName string, maxSize, maxBackups, maxAge int, compress bool, enableKafka bool, kafkaAddress []string) {
	// 打印错误级别的日志
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// 打印所有级别的日志
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	var allCore []zapcore.Core

	hook := lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge,   //days
		Compress:   compress, // disabled by default
	}

	fileWriter := zapcore.AddSync(&hook)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)

	// for human operators.
	consoleEncoder := zapcore.NewConsoleEncoder(NewEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	// kafka
	if len(kafkaAddress) > 0 && enableKafka {
		var (
			kl  LogKafka
			err error
		)
		kl.Topic = "go_framework_log"
		// 设置日志输入到Kafka的配置
		config := sarama.NewConfig()
		//等待服务器所有副本都保存成功后的响应
		config.Producer.RequiredAcks = sarama.WaitForAll
		//随机的分区类型
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true

		kl.Producer, err = sarama.NewSyncProducer(kafkaAddress, config)
		if err != nil {
			fmt.Printf("connect kafka failed: %+v\n", err)
			os.Exit(-1)
		}
		topicErrors := zapcore.AddSync(&kl)
		// 打印在kafka
		kafkaEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
		var kafkaCore zapcore.Core
		if mode == "dev" {
			kafkaCore = zapcore.NewCore(kafkaEncoder, topicErrors, lowPriority)

		} else {
			kafkaCore = zapcore.NewCore(kafkaEncoder, topicErrors, highPriority)

		}
		allCore = append(allCore, kafkaCore)
	}
	if mode == "dev" {
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, fileWriter, lowPriority))
	} else {
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, consoleDebugging, highPriority))
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, fileWriter, highPriority))
	}
	core := zapcore.NewTee(allCore...)

	// From a zapcore.Core, it's easy to construct a Logger.
	Logger = zap.New(core, zap.AddCaller())
	// defer Logger.Sync()
	initlog(fileName, maxSize, maxBackups, maxAge)
	initlogs(mode, fileName, maxSize, maxBackups, maxAge, compress, enableKafka, kafkaAddress)
}

// timeEncoder格式化时间
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",                           // json时时间键
		LevelKey:       "L",                           // json时日志等级键
		NameKey:        "N",                           // json时日志记录器名
		CallerKey:      "C",                           // json时日志文件信息键
		MessageKey:     "M",                           // json时日志消息键
		StacktraceKey:  "S",                           // json时堆栈键
		LineEnding:     zapcore.DefaultLineEnding,     // 友好日志换行符
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 友好日志等级名大小写（info INFO）
		EncodeTime:     TimeEncoder,                   // 友好日志时日期格式化
		EncodeDuration: zapcore.StringDurationEncoder, // 时间序列化
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 日志文件信息（包/文件.go:行号）
	}
}

func initlog(fileName string, maxSize, maxBackups, maxAge int) {
	beego.SetLogFuncCall(true)
	// 打印错误级别的日志
	if beego.BConfig.RunMode == "dev" {
		beego.SetLevel(beego.LevelDebug)
		// orm.Debug = true
	} else {
		beego.SetLevel(beego.LevelInformational)
	}

	cfg := map[string]interface{}{
		"filename": fileName,
		"maxsize":  maxSize * 1024 * 1024,
		"maxdays":  maxAge,
	}
	cfgstr, _ := json.Marshal(cfg)
	beego.SetLogger("file", string(cfgstr))
}

func initlogs(mode string, fileName string, maxSize, maxBackups, maxAge int, compress bool, enableKafka bool, kafkaAddress []string) {

	option := make([]LogOption, 0)
	option = append(option, SetMaxFileSize(maxSize), SetCompress(compress), SetMaxAge(maxAge), SetMaxBackups(maxBackups))
	level := InfoLevel
	if mode == "dev" {
		level = DebugLevel
		option = append(option, SetCaller(true))
	}
	// Path        string // 文件绝对地址，如：/home/homework/neso/file.log
	// Level       string // 日志输出的级别
	// MaxFileSize int    // 日志文件大小的最大值，单位(M)
	// MaxBackups  int    // 最多保留备份数
	// MaxAge      int    // 日志文件保存的时间，单位(天)
	// Compress    bool   // 是否压缩
	// Caller

	Init(mode, fileName, level, false, option...)
}
