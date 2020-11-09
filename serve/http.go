package serve

import (

	// "time"

	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"

	_ "github.com/qiaogw/initial"
	"github.com/qiaogw/pkg/cache"
	conf1 "github.com/qiaogw/pkg/config"

	//"github.com/qiaogw/pkg/conf1"
	"time"

	"github.com/qiaogw/pkg/consts"
	"github.com/qiaogw/pkg/dbtools"
	"github.com/qiaogw/pkg/jwt"
	"github.com/qiaogw/pkg/logs"
	"github.com/qiaogw/pkg/tools"
	_ "github.com/qiaogw/routers"

	"github.com/astaxie/beego"
	"go.uber.org/zap"

	// _ "github.com/lib/pq"

	//	"github.com/beego/i18n"

	// "github.com/astaxie/beego/orm"

	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	// if beego.BConfig.RunMode == "dev" {
	// 	orm.Debug = true
	// }
}

func Start() {

	beego.Debug("http server start...")
	defer func() {
		if r := recover(); r != nil {
			logs.Fatal("recovered panic", zap.String("type", consts.PanicRecoveredError))
			panic(r)
		}
	}()
	logs.Sync()

	Exit := func(code int) {
		delPidFile()
		os.Exit(code)
	}
	tools.CheckPath(conf1.Config.LockFilePath)
	f := tools.LockOrDie(conf1.Config.LockFilePath)
	defer f.Unlock()

	if err := tools.MakeDirectory(conf1.Config.TempDir); err != nil {
		logs.Fatal("can't create temporary directory", zap.String("type", consts.IOError), zap.String("dir", conf1.Config.TempDir))
		Exit(1)
	}

	// killOld()

	rand.Seed(time.Now().UTC().UnixNano())

	// save the current pid and version
	if err := savePid(); err != nil {
		logs.Fatal("can't create pid", zap.Error(err))

		Exit(1)
	}
	// defer delPidFile()
	beegoRun()
}

func beegoRun() {

	cfg := conf1.Config.DB
	dbtools.GetDBConnection(cfg, "")
	cache.InitCache()
	jwt.InitJwt()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		// orm.Debug = true
	}

	beego.SetStaticPath("/", "dist")
	beego.BConfig.WebConfig.StaticDir["/imgs"] = "imgs"
	// beego.BConfig.Log.AccessLogs = true
	// logLevel, _ := strconv.Atoi(beego.AppConfig.String("log_level"))
	// if beego.BConfig.RunMode == "prod" {
	// 	beego.SetLevel(logLevel)
	// }
	corsConfig := conf1.Config.Cors
	corsOptions := new(cors.Options)
	corsOptions.AllowOrigins = corsConfig.AllowOrigins
	corsOptions.AllowCredentials = corsConfig.AllowCredentials
	corsOptions.AllowMethods = corsConfig.AllowMethods
	corsOptions.AllowHeaders = corsConfig.AllowHeaders
	corsOptions.ExposeHeaders = corsConfig.ExposeHeaders
	corsOptions.MaxAge = corsConfig.MaxAge * time.Hour
	// beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"OPTIONS", "HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
	// 	AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "x-requested-with", "X-Requested-With"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(corsOptions))
	//beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}

func delPidFile() {
	// beego.Debug("delPidFile")
	os.Remove(conf1.Config.GetPidPath())
}

func savePid() error {
	pid := os.Getpid()
	PidAndVer, err := json.Marshal(map[string]string{"pid": tools.IntToStr(pid), "version": consts.VERSION})
	if err != nil {
		logs.Error("marshalling pid to json", zap.String("type", consts.JSONMarshallError), zap.Error(err), zap.Int("pid", pid))
		return err
	}

	return ioutil.WriteFile(conf1.Config.GetPidPath(), PidAndVer, 0644)
}
