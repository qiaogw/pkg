package service

import (
	"io"
	"log"
	"os"
)

func newLogger(writer io.Writer) *consoleLogger {
	if writer == nil {
		writer = os.Stderr
	}
	logger := &consoleLogger{
		info: log.New(writer, "[I] ", log.LstdFlags),
		warn: log.New(writer, "[W] ", log.LstdFlags),
		err:  log.New(writer, "[E] ", log.LstdFlags),
	}
	return logger
}

type consoleLogger struct {
	info, warn, err *log.Logger
}

func (c consoleLogger) Error(v ...interface{}) error {
	c.err.Print(v...)
	return nil
}
func (c consoleLogger) Warning(v ...interface{}) error {
	c.warn.Print(v...)
	return nil
}
func (c consoleLogger) Info(v ...interface{}) error {
	c.info.Print(v...)
	return nil
}
func (c consoleLogger) Errorf(format string, a ...interface{}) error {
	c.err.Printf(format, a...)
	return nil
}
func (c consoleLogger) Warningf(format string, a ...interface{}) error {
	c.warn.Printf(format, a...)
	return nil
}
func (c consoleLogger) Infof(format string, a ...interface{}) error {
	c.info.Printf(format, a...)
	return nil
}
