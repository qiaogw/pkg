package s3cli

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/qiaogw/log"
	"github.com/qiaogw/pkg/config"
)

func Mount() {
	var wg sync.WaitGroup
	go s3Run(wg)
	wg.Wait()
	time.Sleep(3 * time.Second)
}

func s3Run(wg sync.WaitGroup) {
	LocalConf := config.Config.LocalStore
	//cmdRclone := fmt.Sprintf("rclone cmount s3:  %s --cache-dir %s --config %s --vfs-cache-mode writes --allow-other   -q ", LocalConf.Dir, LocalConf.CacheDir, LocalConf.ConfigFile)
	//cmdRclone := fmt.Sprintf("rclone rmdirs s3:buck1  --config %s --leave-root -vv ", LocalConf.ConfigFile)
	wg.Add(1)
	defer wg.Done()
	bucket := "s3:" + LocalConf.BucketsPath
	logStr := fmt.Sprintf("--log-file=%s", LocalConf.LogFile)
	command := exec.Command("rclone", "mount", bucket, LocalConf.Dir, "--cache-dir", LocalConf.CacheDir, "--config", LocalConf.ConfigFile, "--vfs-cache-mode", "writes", "--allow-other", "--drive-use-trash=false", logStr, "--allow-non-empty", "-q")
	//start the execution
	if err := command.Start(); err != nil {
		log.Error("Failed to start cmd: ", err)
	}
	if err := command.Wait(); err != nil {
		log.Error("Failed to Wait cmd: ", err) //
	}
}
