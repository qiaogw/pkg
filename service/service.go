package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/kardianos/service"
	"github.com/qiaogw/pkg/logs"
)

func ValidServiceAction(action string) error {
	for _, act := range service.ControlAction {
		if act == action {
			return nil
		}
	}
	return fmt.Errorf("Available actions: %q", service.ControlAction)
}

// New 以服务的方式启动
// 服务支持的操作有：
// emanager service install  	-- 安装服务
// emanager service uninstall  -- 卸载服务
// emanager service start 		-- 启动服务
// emanager service stop 		-- 停止服务
// emanager service restart 	-- 重启服务
func New(cfg *Config, action string) error {
	p := NewProgram(cfg)
	p.logger.Infof("servic %s new  %s ...", p.DisplayName, action)
	sys := service.ChosenSystem()
	s, err := sys.New(p, &p.Config.Config)

	if err != nil {
		p.logger.Error("servic new   ...", p.DisplayName)
		return err
	}
	p.service = s
	// Service
	if action != "run" {
		//if err := ValidServiceAction(action); err != nil {
		//	p.logger.Errorf(" ValidServiceAction servic action  %s err is %v  ...", action, err)
		//	return err
		//}
		err = service.Control(s, action)
		if err != nil {
			p.logger.Errorf("Valid action  %s err is %v  ", action, err)
			logs.Fatal(err)
		}
		return err
	}
	err = s.Run()
	if err != nil {
		p.logger.Errorf("servic run  %s err is %v  ...", p.DisplayName, err)
	}
	return err
}

func getPidFiles(pidFilePath string) []string {
	var pidFile []string
	// pidFilePath := config.Config.GetPidPath()
	err := filepath.Walk(pidFilePath, func(pidPath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if filepath.Ext(pidPath) == `.pid` {
			pidFile = append(pidFile, pidPath)
		}
		return nil
	})
	if err != nil {
		logs.Error(err)
	}
	return pidFile
}

func NewProgram(cfg *Config) *program {
	// pidFile := config.Config.GetPidPath()
	pidFile := cfg.PidFilePath
	err := os.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		logs.Error(err)
	}
	pidName := cfg.Name + ".pid"
	pidFile = filepath.Join(pidFile, pidName)
	p := &program{
		Config:  cfg,
		pidFile: pidFile,
	}
	p.Config.Config.Arguments = append([]string{`service`, `run`}, p.Args...)
	p.Config.Config.WorkingDirectory = p.Dir
	p.Config.Config.Option = map[string]interface{}{
		//`PIDFile`:   pidFile,
		`RunAtLoad`: true,
		//`UserService`:   true, //Install as a current user service.
		//`SessionCreate`: true, //Create a full user session.
	}
	return p
}

type program struct {
	*Config
	service  service.Service
	cmd      *exec.Cmd
	fullExec string
	pidFile  string
}

func (p *program) Start(s service.Service) (err error) {
	p.logger.Info("servic starting   ...", p.DisplayName)
	if service.Interactive() {
		p.logger.Info("Running in terminal.")
	} else {
		p.logger.Info("Running under service manager.")
	}
	if filepath.Base(p.Exec) == p.Exec {
		p.fullExec, err = exec.LookPath(p.Exec)
		if err != nil {
			p.logger.Errorf("Failed to find executable %q: %v", p.Exec, err)
			return fmt.Errorf("Failed to find executable %q: %v", p.Exec, err)
		}
	} else {
		p.fullExec = p.Exec
	}
	p.createCmd()

	go p.run()
	return nil
}

func (p *program) createCmd() {
	if len(p.Args) < 1 {
		p.Args = append(p.Args, "start")
	}
	p.cmd = exec.Command(p.fullExec, p.Args...)

	p.cmd.Dir = p.Dir
	p.cmd.Env = append(os.Environ(), p.Env...)
	if p.Stderr != nil {
		p.cmd.Stderr = p.Stderr
	}
	if p.Stdout != nil {
		p.cmd.Stdout = p.Stdout
	}
	p.logger.Infof("Running cmd: %s %#v", p.fullExec, p.Args)
	p.logger.Infof("Workdir: %s", p.cmd.Dir)
	//p.logger.Infof("Env: %s", com.Dump(p.cmd.Env, false))
}

func (p *program) Stop(s service.Service) error {
	p.logger.Info("servic Stop:", p.Name)
	p.killCmd()
	p.logger.Infof("Stopping %s", p.DisplayName)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *program) killCmd() {
	//err := com.CloseProcessFromCmd(p.cmd)
	//if err != nil {
	//	p.logger.Error(err)
	//}
	//err = com.CloseProcessFromPidFile(p.pidFile)
	//if err != nil {
	//	p.logger.Error(p.pidFile+`:`, err)
	//}
	//for _, pidFile := range getPidFiles(p.PidFilePath) {
	//	err = com.CloseProcessFromPidFile(pidFile)
	//	if err != nil {
	//		p.logger.Error(pidFile+`:`, err)
	//	}
	//}
}

func (p *program) close() {
	if service.Interactive() {
		p.Stop(p.service)
	} else {
		p.service.Stop()
		p.killCmd()
	}
	if p.Config.OnExited != nil {
		err := p.Config.OnExited()
		if err != nil {
			p.logger.Error(err)
		}
	}
}

func (p *program) run() {
	p.logger.Infof("Starting runing %s", p.DisplayName)
	//return
	//如果调用的程序停止了，则本服务同时也停止
	defer p.close()
	err := p.cmd.Start()
	if err == nil {
		p.logger.Info("APP PID:", p.cmd.Process.Pid)
		ioutil.WriteFile(p.pidFile, []byte(strconv.Itoa(p.cmd.Process.Pid)), os.ModePerm)
		//err = p.cmd.Wait()
		if err := p.cmd.Wait(); err != nil {
			p.logger.Error("Failed to Wait cmd: ", err) //
		}
	}
	if err != nil {
		p.logger.Error("Error running:", err)
	}
}
