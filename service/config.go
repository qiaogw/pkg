package service

import (
	"github.com/kardianos/service"
	"io"
	"os"
)

type Options struct {
	Name        string // Required name of the service. No spaces suggested.
	DisplayName string // Display name, spaces allowed.
	Description string // Long description of service.
	LogFile     string // 服务日志.
}

// Config is the runner app config structure.
type Config struct {
	service.Config
	logger service.Logger

	Dir         string
	Exec        string
	Args        []string
	Env         []string
	PidFilePath string

	OnExited       func() error `json:"-"`
	Stderr, Stdout io.Writer    `json:"-"`
}

func (c *Config) CopyFromOptions(options *Options) *Config {
	c.Name = options.Name
	c.DisplayName = options.DisplayName
	c.Description = options.Description
	return c
}

func FileWriter(file string) (io.WriteCloser, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	return f, err
}
