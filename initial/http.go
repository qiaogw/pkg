package initial

//
//import (
//	"github.com/qiaogw/pkg/config"
//	"github.com/qiaogw/pkg/dbtools"
//	"time"
//
//	"github.com/astaxie/beego"
//	"github.com/astaxie/beego/plugins/cors"
//)
//
//type InitController struct {
//	beego.Controller
//}
//
//func (this *InitController) Get() {
//	data := config.Config
//	this.Data["json"] = data
//	this.ServeJSON()
//}
//
//func (this *InitController) Setup() {
//	err := dbtools.DbSetup()
//	data := config.Config.DB
//	if err != nil {
//		beego.Error(err)
//		this.Data["json"] = err.Error()
//	} else {
//		this.Data["json"] = data
//	}
//	// this.Data["json"] = data
//	this.ServeJSON()
//	// this.StopRun()
//
//}
//
//func (this *InitController) Backup() {
//	cfg := config.Config.DB
//	ft, err := dbtools.BackupDB(cfg)
//	if err != nil {
//		beego.Error(err)
//		this.Data["json"] = err
//	} else {
//		this.Data["json"] = ft
//	}
//	// data := conf.Config.DB
//	// this.Data["json"] = data
//	this.ServeJSON()
//}
//
//func Start() {
//	if beego.BConfig.RunMode == "dev" {
//		beego.BConfig.WebConfig.DirectoryIndex = true
//		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
//	}
//	corsConfig := config.Config.Cors
//	corsOptions := new(cors.Options)
//	corsOptions.AllowOrigins = corsConfig.AllowOrigins
//	corsOptions.AllowCredentials = corsConfig.AllowCredentials
//	corsOptions.AllowMethods = corsConfig.AllowMethods
//	corsOptions.AllowHeaders = corsConfig.AllowHeaders
//	corsOptions.ExposeHeaders = corsConfig.ExposeHeaders
//	corsOptions.MaxAge = corsConfig.MaxAge * time.Hour
//	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(corsOptions))
//	// beego.ErrorController(&controllers.ErrorController{})
//
//	beego.Router("/ilive/dbinit", &InitController{})
//	beego.Router("/ilive/setup", &InitController{}, "*:Setup")
//	beego.Router("/ilive/backup", &InitController{}, "*:Backup")
//
//	beego.Run()
//}
