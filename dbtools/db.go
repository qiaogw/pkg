package dbtools

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// "github.com/qiaogw/dadmin/models"
	"io/ioutil"
	"os"
	"path/filepath"

	//"github.com/qiaogw/pkg/conf"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	// _ "github.com/qiaogw/dadmin/models"
	"github.com/qiaogw/pkg/config"
	"github.com/qiaogw/pkg/logs"
	"github.com/qiaogw/pkg/tools"
	// "github.com/xormplus/xorm"
)

var (
	// 初始化安装
	Installed sql.NullBool
	// 安装版本
	installedSchemaVer float64
	// 安装时间
	installedTime time.Time
	// 重现加载
	Orm *orm.Ormer
)

// SetInstalled 设置数据已初始化
func SetInstalled(lockFile string) error {
	now := time.Now()
	err := ioutil.WriteFile(lockFile, []byte(now.Format(`2006-01-02 15:04:05`)+"\n"+fmt.Sprint(config.Config.DBUpdateToVersion)), os.ModePerm)
	if err != nil {
		return err
	}
	installedTime = now
	installedSchemaVer = config.Config.DBUpdateToVersion
	Installed.Valid = true
	Installed.Bool = true
	return nil
}

//IsInstalled 判断数据是否已经初始化
func IsInstalled() bool {
	dir, _ := os.Getwd()
	if !Installed.Valid {
		lockFile := filepath.Join(dir, `installed.lock`)
		if info, err := os.Stat(lockFile); err == nil && !info.IsDir() {
			if b, e := ioutil.ReadFile(lockFile); e == nil {
				content := string(b)
				content = strings.TrimSpace(content)
				lines := strings.Split(content, "\n")
				switch len(lines) {
				case 2:
					installedSchemaVer, _ = strconv.ParseFloat(strings.TrimSpace(lines[1]), 64)
					fallthrough
				case 1:
					installedTime, _ = time.Parse(`2006-01-02 15:04:05`, strings.TrimSpace(lines[0]))
				}
			}
			Installed.Valid = true
			Installed.Bool = true
		}
	}

	if Installed.Bool && config.Config.DBUpdateToVersion > installedSchemaVer {
		var upgraded bool
		// err := createdb()
		// if err == nil {
		// 	upgraded = true
		// } else {
		// 	logs.Panic(`数据库表结构需要升级！`)
		// }
		if upgraded {
			installedSchemaVer = config.Config.DBUpdateToVersion
			ioutil.WriteFile(filepath.Join(dir, `installed.lock`), []byte(installedTime.Format(`2006-01-02 15:04:05`)+"\n"+fmt.Sprint(config.Config.DBUpdateToVersion)), os.ModePerm)
		}
	}
	return Installed.Bool
}

func GetDBConnectionString(cfg config.DBConfig) (str string) {
	switch cfg.DbType {
	case "mysql":
		str = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.Charset)
		orm.RegisterDriver("mysql", orm.DRMySQL)
	case "postgres":
		str = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password)
		orm.RegisterDriver("postgres", orm.DRPostgres)
	case "sqlite3":
		cwd, err := os.Getwd()
		if err != nil {
			return
		}
		if cfg.DbPath == "" {
			cfg.DbPath = filepath.Join(cwd, "DB-data")
			config.Config.DB.DbPath = cfg.DbPath
		}
		str = filepath.Join(cfg.DbPath, cfg.Name+".db")
		if err := tools.CheckPath(str); err != nil {
			return
		}
		orm.RegisterDriver("sqlite3", orm.DRSqlite)
	default:
		beego.Error("Database driver is not allowed:", cfg.DbType)
	}
	return
}

//Createdb 创建数据库
func Createdb(cfg config.DBConfig) (err error) {
	var sqlstring string
	//beego.Debug(cfg)
	dns := GetDBConnectionString(cfg)
	db, err := sqlx.Open(cfg.DbType, dns)
	defer db.Close()
	if err != nil {
		logs.Fatal("创建xorm引擎失败", err)
		return
	}
	switch cfg.DbType {
	case "mysql":
		sqlstring = fmt.Sprintf("CREATE DATABASE  if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci;", cfg.Name)
	case "postgres":
		sqlstring = fmt.Sprintf("CREATE DATABASE %s;", cfg.Name)
	case "sqlite3":
		dns := filepath.Join(cfg.DbPath, cfg.Name+".db")
		os.Remove(dns)
		sqlstring = "create table init (n varchar(32));drop table init;"
	default:
		logs.Panic("Database driver is not allowed:", cfg.DbType)
		return
	}
	_, err = db.Exec(sqlstring)
	if err != nil {
		beego.Error("err is ", err)
	}
	return
}

// DbTest 加载数据
func DbTest(cfg config.DBConfig) (err error) {
	connStr := GetDBConnectionString(cfg)
	beego.Debug("正在链接数据库 %s ... ", connStr, cfg)
	tbname := "test"
	err = orm.RegisterDataBase(tbname, cfg.DbType, connStr)
	if err != nil {
		if !strings.Contains(err.Error(), "DataBase alias name `test` already registered") {
			return
		} else {
			err = orm.RegisterDataBase(tbname+tools.UUID(), cfg.DbType, connStr)
		}
	}
	return
}

// DBConnect is initializes beego connection
func DBConnect(cfg config.DBConfig) (err error) {
	connStr := GetDBConnectionString(cfg)
	beego.Debug("正在链接数据库 %s ... ", connStr)
	err = orm.RegisterDataBase("default", cfg.DbType, connStr)
	if err != nil {
		if !strings.Contains(err.Error(), "DataBase alias name `default` already registered") {
			if err := orm.RegisterDataBase("admin", cfg.DbType, connStr); err != nil {
				return err
			}
			if err := orm.NewOrm().Using("admin"); err != nil {
				return err
			}
		}
	}
	beego.Info("数据库链接成功 %s ... ", connStr)
	return nil
}

// Syncdb 创建数据库表
func Syncdb(cfg config.DBConfig, force bool) (err error) {
	dbName := "default"
	_, err = orm.GetDB(dbName)
	if err != nil {
		dbName = "admin"
	}
	err = orm.RunSyncdb(dbName, force, true)
	if err != nil {
		return err
	}
	return
}
func backupTableData(rname string) {
	o := orm.NewOrm()
	var maps []orm.Params
	o.QueryTable(rname).Limit(-1).Values(&maps)
	data, _ := json.Marshal(maps)
	backupFile := filepath.Join(config.Config.DB.BackupPath, rname+".json")
	ioutil.WriteFile(backupFile, data, 0644)
}

func backup() {
	backupTableData("dept")
	backupTableData("fsm_action")
	backupTableData("fsm_event")
	backupTableData("fsm_node")
	backupTableData("organization")
	backupTableData("operation")
	backupTableData("param")
	backupTableData("resource")
	backupTableData("user")
	backupTableData("role")
	backupTableData("permission")
}

/**
* 初始化数据
 */
// func Initdata2() (err error) {
// 	jdata, _ := ioutil.ReadFile("DB-data/operation.json")
// 	var operation []models.Operation
// 	err = json.Unmarshal(jdata, &operation)
// 	beego.Error(err)
// 	for _, v := range operation {
// 		_, err = models.AddOperation(&v)
// 		beego.Error(err)
// 	}
// 	jdata, _ = ioutil.ReadFile("DB-data/param.json")
// 	var param []models.Param
// 	err = json.Unmarshal(jdata, &param)
// 	beego.Error(err)
// 	for _, v := range param {
// 		_, err = orm.NewOrm().Insert(&v)
// 		beego.Error(err)
// 	}
// 	jdata, _ = ioutil.ReadFile("DB-data/resource.json")
// 	var resource []models.Resource
// 	err = json.Unmarshal(jdata, &resource)
// 	beego.Error(err)
// 	for _, v := range resource {
// 		_, err = models.AddResource(&v)
// 		beego.Error(err)
// 	}
// 	return
// }
