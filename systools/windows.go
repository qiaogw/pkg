// +build windows

// Copyright 2018 cloudy itcloudy@qq.com.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package systools

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/qiaogw/pkg/logs"

	"go.uber.org/zap"
)

// KillPid kills the process with the specified pid
func KillPid(pid string) error {

	rez, err := exec.Command("tasklist", "/fi", "PID eq "+pid).Output()
	if err != nil {
		logs.Logger.Error("Error executing command", zap.Error(err), zap.String("cmd", "tasklist /fi PID eq"+pid))
		return err
	}
	if string(rez) == "" {
		return fmt.Errorf("null")
	}
	logs.Logger.Error("command execution result", zap.Error(err), zap.String("cmd", "tasklist /fi PID eq "+pid))

	if ok, _ := regexp.MatchString(`(?i)PID`, string(rez)); !ok {
		return fmt.Errorf("null")
	}
	return nil
}
