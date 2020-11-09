package pathtool

//path 系统工具
import (
	//	"fmt"
	//	"io/ioutil"
	"os"
	"runtime"

	//	"strings"
	//	"syscall"

	"github.com/astaxie/beego"
)

var ostype = runtime.GOOS

// Diskstatus 磁盘状态
type Diskstatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

//FileSize 获取文件大小
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

//DiskUsage 获取磁盘空间
func DiskUsage(path string) (disk Diskstatus) {
	//	if ostype == "windows" {

	//	} else if ostype == "linux" {
	//		fs := syscall.Statfs_t{}
	//		err := syscall.Statfs(path, &fs)
	//		if err != nil {
	//			beego.Error(err)
	//			return
	//		}
	//		disk.All = fs.Blocks * uint64(fs.Bsize)
	//		disk.Free = fs.Bfree * uint64(fs.Bsize)
	//		disk.Used = disk.All - disk.Free
	//	}

	return
}

//CheckPath 检查目录是否存在，若不存在则创建目录
//输入为路径名称
//返回为错误
func CheckPath(logPath string) error {
	//	dir := filepath.Dir(logPath)
	_, err := os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(logPath, os.ModePerm)
			if err != nil {
				beego.Error("创建文件夹失败！：", err.Error())
				return err
			}
			//			Info("文件夹不存在但已创建！")
		}
	}
	beego.Info(logPath, " 文件夹已存在！")
	return nil
}

// CheckFileIsExist 判断文件是否存在  存在返回 true 不存在返回false
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
