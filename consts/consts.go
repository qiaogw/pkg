package consts

const VERSION = "0.1.0"

//default value
const (
	DefaultAppName = "emanager"
	//user upload file url
	USER_UPLOAD_FILE_URL = "/upload_files/"
	//system file url
	SYSTEM_STATIC_FILE_URL = "/system_statics/"
	// default database type
	DefaultDatabase = "postgres"

	// DefaultConfigFile name of config file (toml format)
	DefaultConfigFile = "app.toml"
	// DefaultLogDirName name of config file (toml format)
	DefaultLogDirName = "logs"

	// DefaultConfigDirName name of working directory
	DefaultConfigDirName = "conf"

	// DefaultPidFilename is default filename of pid file
	DefaultPidFilename = "emanager.pid"

	// DefaultLockFilename is default filename of lock file
	DefaultLockFilename = "emanager.lock"
	//DefaultLogFileName
	DefaultLogFileName = "emanager.log"
	//DefaultCaddyLogFileName
	DefaultCaddyLogFileName = "caddy.log"
	// server file dir
	DefaultSystemDataDirName = "system-data"
	// user file upload file dir
	DefaultUserDataDirName = "user-data"
	// temp file dir
	DefaultTempDirName = "app-temp"
	// only for sqllite3 file dir
	DefaultCaddyfile       = "Caddyfile"
	DefaultDbPath          = "DB-data"
	DefaultWebPath         = "web-data"
	DefaultStore           = "local"
	DefaultCacheDir        = "s3cache"
	DefaultStoreDir        = "s3data"
	DefaultStoreConfigFile = "conf/s3/s3.conf"
	DefaultStoreIgnore     = ".s3dataIgnoreIgnoreIgnore"
)

//context variable
const (
	// login user name

	LoginUserName = "LOGIN_USER_NAME"
	// login user id
	LoginUserID = "LOGIN_USER_ID"
	// login user roles []string
	LoginUserRoleIds   = "LOGIN_USER_ROLE_IDS"
	LoginUserRoleCodes = "LOGIN_USER_ROLE_CODES"

	//login user is admin
	LoginIsAdmin = "LOGIN_IS_ADMIN"
	// token is valid
	TokenValid = "TOKEN_VALID"
)
const (
	DefaultPage = 1
	DefaultSize = 20
	DefaultAuth = "qgw"
)
