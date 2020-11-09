//Package pathtool 路径相关
package pathtool

import (
	//	"fmt"
	//	"io/ioutil"
	//	"os"
	//	"runtime"
	"strings"
	//	"syscall"

	"github.com/astaxie/beego"
)

//SplitPath 去除文件nginx根路径
func SplitPath(filepath string) string {
	nginxroot := beego.AppConfig.String("NginxRoot")

	strs := strings.Split(filepath, nginxroot)
	if len(strs) > 1 {
		return strs[1]
	}
	return ""
}
