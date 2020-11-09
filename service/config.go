package service

import (
	"io"
	stdLog "log"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/pkg/errors"
	"github.com/qiaogw/com"
	"github.com/qiaogw/log"
	"github.com/qiaogw/pkg/config"
)

type Options struct {
	Name        string // Required name of the service. No spaces suggested.
	DisplayName string // Display name, spaces allowed.
	Description string // Long description of service.
}

// Config is the runner app config structure.
type Config struct {
	service.Config
	logger service.Logger

	Dir  string
	Exec string
	Args []string
	Env  []string

	OnExited       func() error `json:"-"`
	Stderr, Stdout io.Writer    `json:"-"`
}

func (c *Config) CopyFromOptions(options *Options) *Config {
	c.Name = options.Name
	c.DisplayName = options.DisplayName
	c.Description = options.Description
	return c
}

func Run(options *Options, action string) error {
	conf := &Config{}
	conf.CopyFromOptions(options)
	conf.Dir = com.SelfDir()
	//conf.Exec = os.Args[0]
	conf.Exec = os.Args[0]
	if len(os.Args) > 3 {
		conf.Args = os.Args[3:]
	}

	logDir := config.Config.Log.FilePath
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0775)
		if err != nil {
			return errors.Wrapf(err, "creating dir %s", logDir)
		}
	}

	fileTarget := log.NewFileTarget()
	fileTarget.FileName = filepath.Join(logDir, `app_{date:20060102}.log`) //按天分割日志
	fileTarget.MaxBytes = 10 * 1024 * 1024
	fileTarget.MaxLevel = log.LevelInfo
	log.SetTarget(fileTarget)

	conf.Stderr = log.Writer(log.LevelError)
	conf.Stdout = log.Writer(log.LevelInfo)

	w, err := FileWriter(filepath.Join(logDir, `service.log`))
	if err != nil {
		return err
	}
	conf.OnExited = func() error {
		if w != nil {
			return w.Close()
		}
		return nil
	}
	stdLog.SetOutput(w)
	stdLog.SetFlags(stdLog.Lshortfile)
	conf.logger = newLogger(w)
	return New(conf, action)
}

func FileWriter(file string) (io.WriteCloser, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	return f, err
}
