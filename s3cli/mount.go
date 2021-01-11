package s3cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/qiaogw/pkg/logs"
	//"sync"
	//"time"
)

// Mount mount桶
func Mount(bucket string, s3conf *S3Config) {

	//var wg sync.WaitGroup
	//go s3Run(wg)
	go s3Run(bucket, s3conf)
	//wg.Wait()
	//time.Sleep(3 * time.Second)
	select {}
}

func s3Run(bucket string, s3conf *S3Config) {
	//bucket := "s3:" + config.Config.S3.Bucket
	command := Gets3Cmd(bucket, s3conf)
	//start the execution
	if err := command.Start(); err != nil {
		logs.Error("Failed to start cmd: ", err)
	}
	if err := command.Wait(); err != nil {
		logs.Error("Failed to Wait cmd: ", err) //
	}
}

// Gets3Cmd 获取mount命令
func Gets3Cmd(bucket string, s3conf *S3Config) *exec.Cmd {
	//cmdRclone := fmt.Sprintf("rclone cmount s3:  %s --cache-dir %s --config %s --vfs-cache-mode writes --allow-other   -q ", LocalConf.Dir, LocalConf.CacheDir, LocalConf.ConfigFile)
	//cmdRclone := fmt.Sprintf("rclone rmdirs s3:buck1  --config %s --leave-root -vv ", LocalConf.ConfigFile)
	//wg.Add(1)
	//defer wg.Done()

	logStr := fmt.Sprintf("--log-file=%s", s3conf.LogFile)
	cwd, _ := os.Getwd()
	exefile := filepath.Join(cwd, "rclone")
	gos := runtime.GOOS
	if gos == "windows" {
		exefile = filepath.Join(cwd, "rclone.exe")
	}

	command := exec.Command(exefile, "mount", bucket, s3conf.MountDir, "--cache-dir", s3conf.CacheDir, "--config", s3conf.MountConfigFile, "--vfs-cache-mode", "writes", "--allow-other", "--drive-use-trash=false", logStr, "-q")
	logs.Debug(command.String())
	return command
}
