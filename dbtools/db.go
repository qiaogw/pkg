package dbtools

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	//"github.com/qiaogw/pkg/conf"
	"strconv"
	"strings"
	"time"

	"github.com/qiaogw/pkg/config"
	"github.com/qiaogw/pkg/consts"
	"github.com/qiaogw/pkg/logs"
	"github.com/qiaogw/pkg/tools"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/qiaogw/models"

	// "github.com/xormplus/xorm"
	"github.com/go-xorm/xorm"
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

// LoadInitData 加载初始化数据
func LoadInitData() {
	if config.Config.Init.Enable {
		cfg := config.Config.DB
		GetDBConnection(cfg, "")
		beego.Info("start load init data")
		if config.Config.Init.API != "" {
			beego.Info("start load init api data")
			// initApi()
			beego.Info("load init api data success")

		}
	} else {
		beego.Error("start load init data failed,config file init enable is false")
		os.Exit(-1)
	}

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

func GetDBConnectionString(cfg config.DBConfig, action string) (str string) {
	// cfg := Config.DB
	if action == "" {
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
				beego.Error(err, "getting current wd")
			}
			if cfg.DbPath == "" {
				cfg.DbPath = filepath.Join(cwd, consts.DefaultDbPath)
			}
			// str = fmt.Sprintf("%s%s.db", cfg.DbPath, cfg.Name)
			str = filepath.Join(cfg.DbPath, cfg.Name+".db")
			orm.RegisterDriver("sqlite3", orm.DRSqlite)
		default:
			beego.Error("Database driver is not allowed:", cfg.DbType)
		}
	} else {
		switch cfg.DbType {
		case "mysql":
			str = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8", cfg.User, cfg.Password, cfg.Host, cfg.Port)
		case "postgres":
			str = fmt.Sprintf("host=%s  user=%s  password=%s  port=%d  sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Port, cfg.SslMode)
		case "sqlite3":
			if cfg.DbPath == "" {
				cfg.DbPath = "./"
			}
			str = filepath.Join(cfg.DbPath, cfg.Name+".db")
		default:
			beego.Error("Database driver is not allowed:", cfg.DbType)
		}
	}

	return
}

// GetDBConnection is initializes sqlx connection
func GetDBConnection(cfg config.DBConfig, action string) error {
	var err error
	connStr := GetDBConnectionString(cfg, action)
	if action != "" {
		config.SqlxDB, _ = sqlx.Open(cfg.DbType, connStr)
	} else {
		beego.Debug("正在链接数据库 %s ... ", connStr)
		err = orm.RegisterDataBase("default", cfg.DbType, connStr)
		//初始化之前采用sqllite3
		if err != nil {
			cfg.DbType = "sqlite3"
			cfg.Host = "localhost"
			connStr = GetDBConnectionString(cfg, action)
			beego.Debug(connStr)
			//cwd, _ := os.Getwd()
			beego.Error("cfg is ", cfg)
			//lockFile := filepath.Join(cwd, `installed.lock`)
			if !IsInstalled() {
				err = initDb(cfg, true)
				beego.Error("err is ", err)
			}
			err = orm.RegisterDataBase("default", cfg.DbType, connStr)
			if err != nil {
				beego.Error("err is ", err)
				beego.Error("can't open connection to DB:", cfg.DbType)
				return err
			}
		}
	}
	beego.Info("数据库链接成功 %s ... ", connStr)
	return nil
}

// DbSetup 数据初始化
func DbSetup() (err error) {
	cwd, _ := os.Getwd()
	cfg := config.Config.DB
	lockFile := filepath.Join(cwd, `installed.lock`)
	if IsInstalled() {
		msg := fmt.Sprintf("已经安装过了。如要重新安装，请先删除%s", lockFile)
		return errors.New(msg)
	} else {
		err = initDb(cfg, false)
		if err != nil {
			beego.Error("err is ", err)
			return
		}
		// 生成锁文件
		beego.Info(`Generated file: `, lockFile)
		SetInstalled(lockFile)
	}
	defer config.SqlxDB.Close()
	return
}

func initDb(cfg config.DBConfig, temp bool) (err error) {
	cwd, _ := os.Getwd()
	lockFile := filepath.Join(cwd, `installed.lock`)
	if IsInstalled() {
		msg := fmt.Sprintf("已经安装过了。如要重新安装，请先删除%s", lockFile)
		return errors.New(msg)
	} else {
		err = Createdb(cfg)
		if err != nil {
			return
		}
		dbName := "default"
		_, err := orm.GetDB(dbName)
		if err == nil {
			dbName = "sync"
		}
		connStr := GetDBConnectionString(cfg, "")
		err = orm.RegisterDataBase(dbName, cfg.DbType, connStr)
		err = orm.RunSyncdb(dbName, temp, true)
		if err != nil {
			return err
		}

		err = InitData(cfg, temp)
		if err != nil {
			return err
		}
	}
	defer config.SqlxDB.Close()
	return
}

//Createdb 创建数据库
func Createdb(cfg config.DBConfig) (err error) {
	var sqlstring string
	//beego.Debug(cfg)
	GetDBConnection(cfg, "init")
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
	}
	db := config.SqlxDB
	if err != nil {
		beego.Error("err is ", err)
	}
	_, err = db.Exec(sqlstring)
	if err != nil {
		beego.Error("err is ", err)
	}
	defer db.Close()
	return
}

// InitData 加载数据
func InitData(cfg config.DBConfig, temp bool) error {
	connStr := GetDBConnectionString(cfg, "")
	fileDir := filepath.Join(config.Config.DB.DbPath, beego.BConfig.AppName, beego.BConfig.AppName+".sql")
	db, _ := xorm.NewEngine(cfg.DbType, connStr)
	_, err := db.ImportFile(fileDir)
	defer db.Close()
	return err
}

func backupData(rname string) (err error) {
	o := orm.NewOrm()
	var lists []orm.ParamsList
	tableName := "public." + rname
	sqlStr := fmt.Sprintf("SELECT * FROM %v ", tableName)
	_, _ = o.Raw(sqlStr).ValuesList(&lists)
	dateStr := tools.GetDateMH(time.Now().Unix())
	filePath := filepath.Join(config.Config.DB.DbPath, beego.BConfig.AppName, beego.BConfig.AppName+"-"+dateStr+".sql")
	filePathNew := filepath.Join(config.Config.DB.DbPath, beego.BConfig.AppName, beego.BConfig.AppName+".sql")
	for _, v := range lists {
		strV := tools.ArrayToStr(v)
		str := fmt.Sprintf(`INSERT INTO %v VALUES ( %v );`, tableName, strV)
		tools.WriteFile(filePath, str, false)
	}
	_, err = tools.CopyFile(filePath, filePathNew)
	return
}

func BackupDBData() {
	o := orm.NewOrm()
	sqlStr := "SELECT tablename FROM pg_tables where schemaname='public'"
	var list orm.ParamsList
	num, err := o.Raw(sqlStr).ValuesFlat(&list)
	if err == nil && num > 0 {
		filePathSqlite := filepath.Join(config.Config.DB.DbPath, beego.BConfig.AppName, beego.BConfig.AppName+"-sqlite"+".sql")
		tools.WriteFile(filePathSqlite, "--###Sqllite.sql##", true)
		for _, tn := range list {
			backupData(tn.(string))
		}
	}
}
