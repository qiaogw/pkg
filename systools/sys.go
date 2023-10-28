package systools

import (

	// "time"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/qiaogw/pkg/logs"

	"time"

	"go.uber.org/zap"
	// _ "github.com/lib/pq"
	//	"github.com/beego/i18n"
	// "github.com/astaxie/beego/orm"
)

func KillOld() {
	pidPath := "path"
	if _, err := os.Stat(pidPath); err == nil {
		dat, err := ioutil.ReadFile(pidPath)
		if err != nil {
			logs.Error("reading pid file failed", zap.String("path", pidPath), zap.Error(err))
		}
		var pidMap map[string]string
		err = json.Unmarshal(dat, &pidMap)
		if err != nil {
			logs.Error("un marshalling pid map", zap.String("type", "JSONUnmarshall"), zap.ByteString("data", dat), zap.Error(err))

		}
		logs.Debug("old pid path", zap.String("path", pidPath))

		KillPid(pidMap["pid"])
		if fmt.Sprintf("%s", err) != "null" {
			// give 15 sec to end the previous process
			for i := 0; i < 15; i++ {
				if _, err := os.Stat(pidPath); err == nil {
					time.Sleep(time.Second)
				} else {
					break
				}
			}
		}
	}
}
