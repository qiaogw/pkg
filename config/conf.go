// Copyright 2018 cloudy itcloudy@qq.com.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
	"github.com/qiaogw/log"
	logConfig "github.com/qiaogw/ozzo-config"
	"github.com/qiaogw/pkg/caddy"
	"github.com/qiaogw/pkg/consts"
	"github.com/qiaogw/pkg/hugo"
	"github.com/qiaogw/pkg/logs"
	"go.uber.org/zap"

	//"context"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// HostPort endpoint in form "str:int"
type HostPort struct {
	Host   string `label:"api地址"` // ipaddr, hostname, or "0.0.0.0"
	Port   int    `label:"api端口"` // must be in range 1..65535
	Enable bool   `label:"api是否有效"`
}

// RpcPort endpoint in form "str:int"
type RpcPort struct {
	Host   string `label:"Rpc地址"` // ipaddr, hostname, or "0.0.0.0"
	Port   int    `label:"Rpc端口"` // must be in range 1..65535
	Enable bool   `label:"Rpc是否启用"`
}

// Str converts HostPort pair to string format
func (h HostPort) Str() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

// DBConfig database connection parameters
type DBConfig struct {
	DbType      string `label:"数据库类型"`   //mysql postgres sqlite3
	Charset     string `label:"数据库编码"`   // only for mysql
	Name        string `label:"数据库名称"`   // db name
	Host        string `label:"数据库地址"`   // ipaddr, hostname, or "0.0.0.0"
	Port        int    `label:"数据库端口"`   // must be in range 1..65535
	User        string `label:"数据库用户名"`  // db user
	Password    string `label:"数据库密码"`   //db password
	LockTimeout int    `label:"事务锁超时时间"` // lock_timeout in milliseconds
	DbPath      string `label:"数据库地址"`   //only for sqlite3
	SslMode     string `label:"SslMode"`
	BackupPath  string `label:"数据库备份地址"`
}

// CentrifugoConfig connection params
type CentrifugoConfig struct {
	Enable bool
	Secret string
	URL    string
}

func (c CentrifugoConfig) String() string {
	return fmt.Sprintf("Secret: %s URL: %s", c.Secret, c.URL)
}

// Log represents parameters of log
type LogConfig struct {
	EnableKafka  bool     `label:"启用kafaka"`
	KafkaAddress []string `label:"kafaka地址"`
	FilePath     string   `label:"文件地址"`       // file name
	MaxSize      int      `label:"日志文件最大尺寸"`   // file size
	MaxBackups   int      `label:"日志文件最大保存"`   // file back
	MaxAge       int      `label:"日志文件最大保存天数"` // file save days
	Compress     bool     `label:"是否压缩"`       // compress file
}

// EmailNotificationConfig smtp config
type EmailNotificationConfig struct {
	Enable   bool   `label:"是否启用"`
	Host     string `label:"smtp地址"`
	Port     int    `label:"端口"`
	Username string `label:"用户名"`
	Password string `label:"密码"`
	To       string `label:"发送至"`
	From     string `label:"发送者"`
	Subject  string `label:"主题"`
}
type TlsConfig struct {
	Enable   bool   `label:"是否启用"`
	CertFile string `label:"CertFile"`
	KeyFile  string `label:"KeyFile"`
}
type RedisConfig struct {
	Enable      bool          `label:"是否启用"`
	Key         string        `label:"Redis collection 的名称"`
	Addr        string        `label:"地址"`
	Password    string        `label:"密码"`
	DBNum       int           `label:"数据库"`
	MaxActive   int           `label:"最大连接数"`
	MaxIdle     int           `label:"最大空闲连接数"`
	IdleTimeout time.Duration `label:"空闲连接超时时间"`
	Wait        bool          `label:"如果超过最大连接，是报错，还是等待。"`
}

type SuperUser struct {
	UserName string `label:"账号"`
	Password string `label:"密码"`
	Email    string `label:"email"`
	Mobile   string `label:"电话"`
}
type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	AllowWebSockets  bool
	MaxAge           time.Duration
}
type InitConfig struct {
	Enable bool   // 是否初始化
	API    string // 接口初始化路径
}
type FileUploadConfig struct {
	Target string `label:"文件上传"` // local,qiniu
}
type ElasticConfig struct {
	Enable        bool
	URLs          []string
	Sniff         bool
	Sniffer       time.Duration
	HealthCheck   bool
	HealthChecker time.Duration
	AuthUserName  string
	AuthPassword  string
	Gzip          bool
}
type Auth struct {
	AuthType int `label:"认证类型"` // 认证类型 0 简单认证 1 登录认证 2 实时认证
}
type Cache struct {
	CacheExpire int    `label:"存活时长"` // cache存活时长，默认60秒
	CacheType   string `label:"类型"`   // file、memcache、memory 和 redis
}
type Seaweed struct {
	SeaweedAddr string `label:"地址"`       // 地址
	Master      string `label:"master地址"` // 地址
	Public      string `label:"Filer只读地址"`
	Private     string `label:"Filer管理地址"`
}
type S3 struct {
	Endpoint        string `label:"地址"`              // 地址
	AccessKeyID     string `label:"AccessKeyID"`     // 地址
	SecretAccessKey string `label:"SecretAccessKey"` // 地址
	Region          string `label:"对象存储的region"`     // 对象存储的region
	Bucket          string `label:"对象存储的Bucket"`     // 对象存储的Bucket
	Secure          bool   `label:"true代表使用HTTPS"`   // true代表使用HTTPS
	Ignore          string `label:"隐藏文件，S3不支持空目录"`   // 地址
	LifeDay         int64  `label:"存储周期，天"`          // 地址
}

type LocalStore struct {
	Dir         string `label:"地址"`            // 地址
	CacheDir    string `label:"Cache地址"`       // 地址
	ConfigFile  string `label:"ConfigFile地址"`  // 地址
	LogFile     string `label:"LogFile地址"`     // 地址
	BucketsPath string `label:"Buckets地址"`     // 地址
	Ignore      string `label:"隐藏文件，S3不支持空目录"` // 地址
	Temp        string `label:"临时目录"`          // 地址
}
type GlobalConfig struct {
	AppName           string                  `label:"应用名称"`
	DBUpdateToVersion float64                 `label:"数据库版本"`            // 数据库升级到某个版本
	RunMode           string                  `label:"运行模式"`             // 运行模式
	TokenExpire       int64                   `label:"token有效时间 秒"`      // token有效时间
	ConfigPath        string                  `label:"配置文件地址"`           // 配置文件地址
	JwtPrivatePath    string                  `label:"jwt private path"` // jwt private path
	JwtPublicPath     string                  `label:"jwt public path"`  // jwt public path
	JwtFlash          bool                    `label:"jwt 是否自动刷新"`
	TempDir           string                  `label:"临时文件"`   // 临时文件
	SystemDataDir     string                  `label:"系统文件地址"` // 系统文件地址，系统自带文件
	UserDataDir       string                  `label:"用户文件地址"` // 用户文件地址，用户上传文件存储位置
	PidFilePath       string                  `label:"进程"`     // 进程
	LockFilePath      string                  // 进程锁 file path
	Store             string                  // 存储配置
	TLS               TlsConfig               // tls
	Redis             RedisConfig             // redis
	HTTP              HostPort                // http端口
	RPC               RpcPort                 // rpc端口
	DB                DBConfig                // 数据库配置
	Centrifugo        CentrifugoConfig        // Centrifugo消息推送
	Log               LogConfig               // 日志
	EmailNotification EmailNotificationConfig // 邮件配置
	Cors              CorsConfig              // 跨域配置
	Admin             SuperUser               // 超级用户
	Init              InitConfig              // 初始化数据
	FileUpload        FileUploadConfig        // 文件上传、
	Elastic           ElasticConfig           // elastic配置
	Auth              Auth                    // 认证配置
	Cache             Cache                   // Cache配置
	Seaweed           Seaweed                 // Seaweed配置
	Caddy             caddy.Config            // caddy 配置
	Hugo              hugo.Config             // hugo 配置
	S3                S3                      // S3 配置
	LocalStore        LocalStore              // LocalStore 配置

}

// GetPidPath returns path to pid file
func (c *GlobalConfig) GetPidPath() string {
	return c.PidFilePath
}
func DefaultConfigPath() string {
	p, err := os.Getwd()
	if err != nil {
		logs.Logger.Fatal("getting current wd")
	}
	//return filepath.Join(p, "config", "config.toml")
	return filepath.Join(p, consts.DefaultConfigDirName, "app.toml")
}

// LoadConfig from configFile
// the function has side effect updating global var Config
func LoadConfig(path string) error {
	//log.WithFields(log.Fields{"path": path}).Info("Loading config")
	return LoadConfigToVar(path, &Config)
}
func LoadConfigToVar(path string, v *GlobalConfig) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		InitConfigFile()
		return errors.Errorf("Unable to load config file %s", path)
	}

	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "reading config")
	}
	err = viper.Unmarshal(v)
	if err != nil {
		return errors.Wrapf(err, "marshalling config to global struct variable")
	}
	v.RunMode = beego.BConfig.RunMode
	v.AppName = beego.BConfig.AppName

	logConf := Config.Log
	logFile := filepath.Join(logConf.FilePath, consts.DefaultLogFileName)
	logs.InitLogger(
		Config.RunMode,
		logFile,
		logConf.MaxSize,
		logConf.MaxBackups,
		logConf.MaxAge,
		logConf.Compress,
		logConf.EnableKafka,
		logConf.KafkaAddress)

	logDir := Config.Log.FilePath
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0775)
		if err != nil {
			return errors.Wrapf(err, "creating dir %s", logDir)
		}
	}
	logConfFile := filepath.Join(Config.ConfigPath, "log.json")
	if _, err := os.Stat(logConfFile); err == nil {
		c := logConfig.New()
		c.Load("app.json")
		// 注册供 Logger.Targets 使用的日志标的类型
		c.Register("ConsoleTarget", log.NewConsoleTarget)
		c.Register("FileTarget", log.NewFileTarget)
	}

	fileTarget := log.NewFileTarget()
	fileTarget.FileName = filepath.Join(logDir, `app_{date:20060102}.log`) //按天分割日志
	fileTarget.MaxBytes = 10 * 1024 * 1024
	fileTarget.MaxLevel = log.LevelInfo
	log.SetTarget(fileTarget)

	log.DefaultLog.Logger = log.GetLogger(Config.AppName, func(l *log.Logger, e *log.Entry) string {
		return fmt.Sprintf("%v [%v] [%v] %v %v", e.Time.Format(time.RFC3339), e.Level, e.Category, e.Message, e.CallStack)
	})

	redis, _ := json.Marshal(viper.Get("RedisProd"))
	db, _ := json.Marshal(viper.Get("DBProd"))
	seaweed, _ := json.Marshal(viper.Get("SeaweedProd"))
	if v.RunMode == "dev" {
		redis, _ = json.Marshal(viper.Get("RedisDev"))
		db, _ = json.Marshal(viper.Get("DBDev"))
		seaweed, _ = json.Marshal(viper.Get("SeaweedDev"))
	}
	prd := new(RedisConfig)
	pdb := new(DBConfig)
	psw := new(Seaweed)
	err = json.Unmarshal(redis, prd)
	if err != nil {
		return errors.Wrapf(err, "marshalling config to global struct variable")
	}
	v.Redis = *prd
	err = json.Unmarshal(db, pdb)
	if err != nil {
		return errors.Wrapf(err, "marshalling config to global struct variable")
	}
	v.DB = *pdb
	err = json.Unmarshal(seaweed, psw)
	if err != nil {
		return errors.Wrapf(err, "marshalling config to global struct variable")
	}
	v.Seaweed = *psw
	return nil
}

func InitConfigFile() error {
	err := FillRuntimePaths()
	if err != nil {
		fmt.Printf("Filling config: %+v\n", err)

		logs.Fatalf("Filling config", zap.Error(err))
	}
	configPath := filepath.Join("conf", consts.DefaultConfigFile)
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Errorf("Marshalling config to global struct variable: %+v\n", err)
	}
	err = SaveConfig(configPath)
	if err != nil {
		fmt.Printf("Saving config failed: %+v\n", err)
	}
	fmt.Println("config file is saved success and path is " + configPath)
	return err
}

// GetConfigFromPath read config from path and returns GlobalConfig struct
func GetConfigFromPath(path string) (interface{}, error) {
	//log.WithFields(log.Fields{"path": path}).Info("Loading config")

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, errors.Errorf("Unable to load config file %s", path)
	}

	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "reading config")
	}

	c := new(interface{})
	err = viper.Unmarshal(c)
	if err != nil {
		return c, errors.Wrapf(err, "marshalling config to global struct variable")
	}

	return c, nil
}

// SaveConfig save global parameters to configFile
func SaveConfig(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0775)
		if err != nil {
			return errors.Wrapf(err, "creating dir %s", dir)
		}
	}

	cf, err := os.Create(path)
	if err != nil {
		//log.WithFields(log.Fields{"type": consts.IOError, "error": err}).Error("Create config file failed")
		return err
	}
	defer cf.Close()

	err = toml.NewEncoder(cf).Encode(Config)
	if err != nil {
		return err
	}
	return nil
}

// FillRuntimePaths fills paths from runtime parameters
func FillRuntimePaths() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrapf(err, "getting current wd")
	}
	if Config.SystemDataDir == "" {
		Config.SystemDataDir = filepath.Join(cwd, consts.DefaultSystemDataDirName)
	}
	if Config.UserDataDir == "" {
		Config.UserDataDir = filepath.Join(cwd, consts.DefaultUserDataDirName)
	}

	if Config.TempDir == "" {
		Config.TempDir = filepath.Join(cwd, consts.DefaultTempDirName)
	}

	if Config.PidFilePath == "" {
		Config.PidFilePath = filepath.Join(Config.SystemDataDir, consts.DefaultPidFilename)
	}
	if Config.LockFilePath == "" {
		Config.LockFilePath = filepath.Join(Config.SystemDataDir, consts.DefaultLockFilename)
	}
	//beego.Debug("Config.Log.FileName", Config.Log.FileName)
	if Config.Log.FilePath == "" {
		Config.Log.FilePath = filepath.Join(cwd, "logs")
	}
	if Config.JwtPrivatePath == "" {
		Config.JwtPrivatePath = filepath.Join(cwd, "config", "jwt", "tm.rsa")
	}
	if Config.JwtPublicPath == "" {
		Config.JwtPublicPath = filepath.Join(cwd, "config", "jwt", "tm.rsa.pub")
	}
	//https
	if Config.TLS.CertFile == "" {
		Config.TLS.CertFile = filepath.Join(cwd, "config", "https", "cert.pem")
	}
	if Config.TLS.KeyFile == "" {
		Config.TLS.KeyFile = filepath.Join(cwd, "config", "https", "key.pem")
	}

	if Config.Init.API == "" {
		Config.Init.API = filepath.Join(cwd, "init", "api_data.yml")

	}

	if Config.DB.DbPath == "" {
		Config.DB.DbPath = filepath.Join(cwd, consts.DefaultDbPath)
	}
	if Config.DB.BackupPath == "" {
		Config.DB.BackupPath = filepath.Join(cwd, consts.DefaultDbPath)
	}
	if Config.Caddy.Caddyfile == "" {
		Config.Caddy.Caddyfile = filepath.Join(cwd, consts.DefaultConfigDirName, consts.DefaultCaddyfile)
	}
	if Config.Caddy.LogFile == "" {
		Config.Caddy.LogFile = filepath.Join(cwd, consts.DefaultLogDirName, consts.DefaultCaddyLogFileName)
	}
	if Config.Hugo.Dir == "" {
		Config.Hugo.Dir = filepath.Join(cwd, consts.DefaultWebPath)
	}
	return nil
}
func MustOK(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
