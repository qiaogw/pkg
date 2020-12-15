package caddy

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/qiaogw/pkg/tools"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/caddyserver/caddy"
	_ "github.com/caddyserver/caddy/caddyhttp"
	"github.com/caddyserver/caddy/caddytls"
	"github.com/caddyserver/certmagic"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	DefaultConfig = &Config{
		Agreed:                  false,
		CAUrl:                   certmagic.DefaultACME.CA,
		CATimeout:               int64(certmagic.HTTPTimeout),
		DisableHTTPChallenge:    certmagic.DefaultACME.DisableHTTPChallenge,
		DisableTLSALPNChallenge: certmagic.DefaultACME.DisableTLSALPNChallenge,
		ServerType:              `http`,
		CPU:                     `100%`,
		PidFile:                 `./system-data/caddy.pid`,
	}
	DefaultVersion = `2.0.0`
	EnableReload   = true
)

func TrapSignals() {
	caddy.TrapSignals()
}

func (c *Config) fixed(appName string)(err error) {
	if len(c.CAUrl) == 0 {
		c.CAUrl = DefaultConfig.CAUrl
	}
	if c.CATimeout == 0 {
		c.CATimeout = DefaultConfig.CATimeout
	}
	if len(c.ServerType) == 0 {
		c.ServerType = DefaultConfig.ServerType
	}
	if len(c.CPU) == 0 {
		c.CPU = DefaultConfig.CPU
	}
	//path, _ := os.Getwd()
	//pidFile := filepath.Join(path, consts.DefaultSystemDataDirName)
	//err := os.MkdirAll(pidFile, os.ModePerm)
	//if err != nil {
	//	log.Println(err)
	//}
	//c.Caddyfile = filepath.Join(path, consts.DefaultConfigDirName, consts.DefaultCaddyfile)
	//pidFile = filepath.Join(pidFile, `caddy.pid`)
	//c.PidFile = pidFile
	//if len(c.LogFile) == 0 {
	//	logFile := filepath.Join(path, consts.DefaultLogDirName)
	//	err := os.MkdirAll(logFile, os.ModePerm)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	c.LogFile = filepath.Join(logFile, consts.DefaultCaddyLogFileName)
	//} else {
	//	err := os.MkdirAll(filepath.Dir(c.LogFile), os.ModePerm)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}
	err=tools.CheckPath(c.LogFile)
	if err != nil {
		return
	}
	tools.CheckPath(c.PidFile)
	if err != nil {
		return
	}
	c.appName = appName
	c.appVersion = DefaultVersion
	c.Agreed = true
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return
}

type Config struct {
	Agreed                  bool   `json:"agreed"` //Agree to the CA's Subscriber Agreement
	CAUrl                   string `json:"caURL"`  //URL to certificate authority's ACME server directory
	DisableHTTPChallenge    bool   `json:"disableHTTPChallenge"`
	DisableTLSALPNChallenge bool   `json:"disableTLSALPNChallenge"`
	Caddyfile               string `json:"caddyFile"`  //Caddyfile to load (default caddy.DefaultConfigFile)
	CPU                     string `json:"cpu"`        //CPU cap
	CAEmail                 string `json:"caEmail"`    //Default ACME CA account email address
	CATimeout               int64  `json:"caTimeout"`  //Default ACME CA HTTP timeout
	LogFile                 string `json:"logFile"`    //Process log file
	PidFile                 string `json:"-"`          //Path to write pid file
	Quiet                   bool   `json:"quiet"`      //Quiet mode (no initialization output)
	Revoke                  string `json:"revoke"`     //Hostname for which to revoke the certificate
	ServerType              string `json:"serverType"` //Type of server to run

	//---
	EnvFile string `json:"envFile"` //Path to file with environment variables to load in KEY=VALUE format
	Plugins bool   `json:"plugins"` //List installed plugins
	Version bool   `json:"version"` //Show version

	//---
	appVersion string
	appName    string
	instance   *caddy.Instance
	stopped    bool
	ctx        context.Context
	cancel     context.CancelFunc
}

func now() string {
	return time.Now().Format(`2006-01-02 15:04:05`)
}

func (c *Config) Start() error {
	caddy.AppName = c.appName
	caddy.AppVersion = c.appVersion
	certmagic.UserAgent = c.appName + "/" + c.appVersion
	c.stopped = false

	// Executes Startup events
	caddy.EmitEvent(caddy.StartupEvent, nil)

	// Get Caddyfile input
	caddyfile, err := caddy.LoadCaddyfile(c.ServerType)
	if err != nil {
		return err
	}

	if EnableReload {
		c.watchingSignal()
	}

	// Start your engines
	c.instance, err = caddy.Start(caddyfile)
	if err != nil {
		return err
	}
	log.Println(`Caddy`, `Server has been successfully started at `+now())

	// Twiddle your thumbs
	c.instance.Wait()
	return nil
}

// Listen to keypress of "return" and restart the app automatically
func (c *Config) watchingSignal() {
	debug := false
	go func() {
		if debug {
			log.Println(`Caddy`, `listen return ==> `+now())
		}
		in := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if debug {
					log.Println(`Caddy`, `reading ==> `+now())
				}
				var Lf byte = '\n'
				var StrLF string = "\n"
				var StrCRLF string = "\r\n"
				input, _ := in.ReadString(Lf)
				if input == StrLF || input == StrCRLF {
					if debug {
						log.Println(`Caddy`, `restart ==> `+now())
					}
					c.Restart()
				} else {
					if debug {
						log.Println(`Caddy`, `waiting ==> `+now())
					}
				}
			}
		}
	}()
}

func (c *Config) Restart() error {
	defer log.Println(`Caddy`, `Server has been successfully reloaded at `+now())
	if c.instance == nil {
		return nil
	}
	// Get Caddyfile input
	caddyfile, err := caddy.LoadCaddyfile(c.ServerType)
	if err != nil {
		return err
	}
	c.instance, err = c.instance.Restart(caddyfile)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Stop() error {
	c.stopped = true
	defer func() {
		c.cancel()
		log.Println(`Caddy`, `Server has been successfully stopped at `+now())
	}()
	if c.instance == nil {
		return nil
	}
	return c.instance.Stop()
}
func (c *Config) Status() bool {
	return c.stopped
}

func (c *Config) Init(appName string) (err error){
	err=c.fixed(appName)
	certmagic.DefaultACME.Agreed = c.Agreed
	certmagic.DefaultACME.CA = c.CAUrl
	certmagic.DefaultACME.DisableHTTPChallenge = c.DisableHTTPChallenge
	certmagic.DefaultACME.DisableTLSALPNChallenge = c.DisableTLSALPNChallenge
	certmagic.DefaultACME.Email = c.CAEmail
	certmagic.HTTPTimeout = time.Duration(c.CATimeout)
	caddy.PidFile = c.PidFile
	caddy.Quiet = c.Quiet
	caddy.RegisterCaddyfileLoader("flag", caddy.LoaderFunc(c.confLoader))
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(c.defaultLoader))
	log.SetFlags(log.Ldate|log.Lshortfile)
	// Set up process log before anything bad happens
	switch c.LogFile {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "":
		log.SetOutput(ioutil.Discard)
	default:
		log.SetOutput(&lumberjack.Logger{
			Filename:   c.LogFile,
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 10,
		})
	}

	//Load all additional envs as soon as possible
	if err := LoadEnvFromFile(c.EnvFile); err != nil {
		mustLogFatalf("%v", err)
	}

	// Check for one-time actions
	if len(c.Revoke) > 0 {
		err := caddytls.Revoke(c.Revoke)
		if err != nil {
			mustLogFatalf(err.Error())
		}
		fmt.Printf("Revoked certificate for %s\n", c.Revoke)
		os.Exit(0)
	}
	if c.Version {
		fmt.Printf("%s %s\n", c.appName, c.appVersion)
		os.Exit(0)
	}
	if c.Plugins {
		fmt.Println(caddy.DescribePlugins())
		os.Exit(0)
	}
	// Set CPU cap
	err = setCPU(c.CPU)
	if err != nil {
		mustLogFatalf(err.Error())
	}
	return
}

// confLoader loads the Caddyfile using the -conf flag.
func (c *Config) confLoader(serverType string) (caddy.Input, error) {
	if c.Caddyfile == "" {
		return nil, nil
	}
	if c.Caddyfile == "stdin" {
		return caddy.CaddyfileFromPipe(os.Stdin, serverType)
	}
	var contents []byte
	if strings.Contains(c.Caddyfile, "*") {
		// Let caddyfile.doImport logic handle the globbed path
		contents = []byte("import " + c.Caddyfile)
	} else {
		var err error
		contents, err = ioutil.ReadFile(c.Caddyfile)
		if err != nil {
			return nil, err
		}
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       c.Caddyfile,
		ServerTypeName: serverType,
	}, nil
}

// defaultLoader loads the Caddyfile from the current working directory.
func (c *Config) defaultLoader(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(caddy.DefaultConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}

// mustLogFatalf wraps log.Fatalf() in a way that ensures the
// output is always printed to stderr so the user can see it
// if the user is still there, even if the process log was not
// enabled. If this process is an upgrade, however, and the user
// might not be there anymore, this just logs to the process
// log and exits.
func mustLogFatalf(format string, args ...interface{}) {
	if !caddy.IsUpgrade() {
		log.SetOutput(os.Stderr)
	}
	log.Fatalf(format, args...)
}

// setCPU parses string cpu and sets GOMAXPROCS
// according to its value. It accepts either
// a number (e.g. 3) or a percent (e.g. 50%).
func setCPU(cpu string) error {
	var numCPU int

	availCPU := runtime.NumCPU()

	if strings.HasSuffix(cpu, "%") {
		// Percent
		var percent float32
		pctStr := cpu[:len(cpu)-1]
		pctInt, err := strconv.Atoi(pctStr)
		if err != nil || pctInt < 1 || pctInt > 100 {
			return errors.New("invalid CPU value: percentage must be between 1-100")
		}
		percent = float32(pctInt) / 100
		numCPU = int(float32(availCPU) * percent)
		if numCPU < 1 {
			numCPU = 1
		}
	} else {
		// Number
		num, err := strconv.Atoi(cpu)
		if err != nil || num < 1 {
			return errors.New("invalid CPU value: provide a number or percent greater than 0")
		}
		numCPU = num
	}

	if numCPU > availCPU {
		numCPU = availCPU
	}

	runtime.GOMAXPROCS(numCPU)
	return nil
}

// LoadEnvFromFile loads additional envs if file provided and exists
// Envs in file should be in KEY=VALUE format
func LoadEnvFromFile(envFile string) error {
	if envFile == "" {
		return nil
	}

	file, err := os.Open(envFile)
	if err != nil {
		return err
	}
	defer file.Close()

	envMap, err := ParseEnvFile(file)
	if err != nil {
		return err
	}

	for k, v := range envMap {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}

	return nil
}

// ParseEnvFile implements parse logic for environment files
func ParseEnvFile(envInput io.Reader) (map[string]string, error) {
	envMap := make(map[string]string)

	scanner := bufio.NewScanner(envInput)
	var line string
	lineNumber := 0

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		lineNumber++

		// skip lines starting with comment
		if strings.HasPrefix(line, "#") {
			continue
		}

		// skip empty line
		if len(line) == 0 {
			continue
		}

		fields := strings.SplitN(line, "=", 2)
		if len(fields) != 2 {
			return nil, fmt.Errorf("Can't parse line %d; line should be in KEY=VALUE format", lineNumber)
		}

		if strings.Contains(fields[0], " ") {
			return nil, fmt.Errorf("Can't parse line %d; KEY contains whitespace", lineNumber)
		}

		key := fields[0]
		val := fields[1]

		if key == "" {
			return nil, fmt.Errorf("Can't parse line %d; KEY can't be empty string", lineNumber)
		}
		envMap[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}
