package tools

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/theckman/go-flock"
)

// CheckPath 检查文件夹是否存在，不存在则创建
func CheckPath(logPath string) error {
	// 获取 logPath 的父目录路径
	dir := filepath.Dir(logPath)

	// 检查父目录是否存在
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// 父目录不存在，创建该父目录及其所有父级目录
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetFileSuffix 获取文件名中的扩展名
func GetFileSuffix(s string) string {
	fileSuffix := path.Ext(s)
	filenameOnly := strings.TrimSuffix(s, fileSuffix)
	return filenameOnly
}

// LockOrDie 加锁文件夹
func LockOrDie(dir string) (*flock.Flock, error) {
	f := flock.New(dir)
	success, err := f.TryLock()
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, err
	}
	return f, nil
}

// MakeDirectory 创建目录，如果不存在的话
func MakeDirectory(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(dir, 0775)
		}
		return err
	}
	return nil
}

// PathExists 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetAllFile 获取指定路径下特定扩展名的文件名列表
func GetAllFile(pathname string, suffix string) (fileSlice []string) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return
	}
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		} else {
			if suffix != "" {
				if strings.HasSuffix(fi.Name(), suffix) {
					fileSlice = append(fileSlice, fi.Name())
				}
			} else {
				fileSlice = append(fileSlice, fi.Name())
			}
		}
	}
	return
}

// CopyDir 复制文件夹及其中的文件
func CopyDir(srcPath string, destPath string) error {
	// 检测源目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath 不是一个正确的目录！")
			return e
		}
	}
	// 检查目标目录是否存在，不存在则创建
	if destInfo, err := os.Stat(destPath); err != nil {
		if err := MakeDirectory(destPath); err != nil {
			return err
		}
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo 不是一个正确的目录！")
			return e
		}
	}
	// 遍历源目录，并逐个拷贝文件
	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			// 将路径中的反斜杠替换为斜杠
			path := strings.Replace(path, "\\", "/", -1)
			// 根据源路径和目标路径生成新的目标文件路径
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			CopyFile(path, destNewPath)
		}
		return nil
	})
	return err
}

// CopyFile 复制文件
func CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	// 分割目标路径中的目录
	destSplitPathDirs := strings.Split(dest, "/")

	// 检查每级目录是否存在，不存在则创建
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := PathExists(destSplitPath)
			if !b {
				// 创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					return w, err
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

// ReadFile 读取文本文件的内容
func ReadFile(path string) (str string, err error) {
	// 打开文件路径
	fi, err := os.Open(path)
	if err != nil {
		err = errors.New(path + " 不是一个正确的路径！")
		return
	}
	defer fi.Close()
	// 读取文件内容
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		err = errors.New(path + " 读取文件失败！")
		return
	}
	str = string(fd)
	return str, err
}

// WriteFile 写入文本文件的内容
// coverType 为 true 时覆盖写入，为 false 时追加写入
func WriteFile(path, info string, coverType bool) (err error) {
	var fl *os.File
	flag := os.O_WRONLY
	if coverType {
		flag = os.O_TRUNC | os.O_WRONLY
	} else {
		flag = os.O_APPEND | os.O_WRONLY
	}
	if CheckFileIsExist(path) { // 如果文件存在
		fl, err = os.OpenFile(path, flag, os.ModePerm) // 打开文件
	} else {
		fl, err = os.Create(path) // 创建文件
	}
	if err != nil {
		err = errors.New(path + " 打开文件失败！")
		return
	}
	defer fl.Close()
	n, err := fl.WriteString(info + "\n")
	if err == nil && n < len(info) {
		err = errors.New(path + " 写入失败！")
	}
	return
}

// WriteFileByte 写入字节数据到文件
// coverType 为 true 时覆盖写入，为 false 时追加写入
func WriteFileByte(path string, info []byte, coverType bool) (err error) {
	var fl *os.File
	flag := os.O_WRONLY
	if coverType {
		flag = os.O_TRUNC | os.O_WRONLY
	} else {
		flag = os.O_APPEND | os.O_WRONLY
	}
	if CheckFileIsExist(path) { // 如果文件存在
		fl, err = os.OpenFile(path, flag, os.ModePerm) // 打开文件
	} else {
		fl, err = os.Create(path) // 创建文件
	}
	defer fl.Close()
	if err != nil {
		// err = errors.New(path + " 打开文件失败！")
		return
	}
	n, err := fl.Write(info)
	if err == nil && n < len(info) {
		err = errors.New(path + " 写入失败！")
	}
	return
}

// CheckFileIsExist 判断文件是否存在，存在返回 true，不存在返回 false
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// GetDirList 获取指定目录下的目录结构树
func GetDirList(dirpath, pathStr string) []DirBody {
	var allFile []DirBody
	finfo, _ := ioutil.ReadDir(dirpath)
	for _, x := range finfo {
		var chiledren DirBody
		if x.IsDir() {
			realPath := filepath.Join(dirpath, x.Name())
			realDir := filepath.Join(pathStr, x.Name())
			chiledren.Label = x.Name()
			chiledren.Dir = realDir
			chiledren.Icon = "el-icon-folder"
			chiledren.Children = append(chiledren.Children, GetDirList(realPath, realDir)...)
			allFile = append(allFile, chiledren)
		}
	}
	return allFile
}

// DirBody 用于表示目录结构的结构体
type DirBody struct {
	Label    string    `json:"label"`
	Children []DirBody `json:"children"`
	Icon     string    `json:"icon"`
	Dir      string    `json:"dir"`
}

// FindEmptyFolder 递归查找空目录
func FindEmptyFolder(dirname string) (emptys []string, err error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	// 判断目录是否为空
	if len(files) == 0 {
		return []string{dirname}, nil
	}

	for _, file := range files {
		if file.IsDir() {
			edirs, err := FindEmptyFolder(path.Join(dirname, file.Name()))
			if err != nil {
				return nil, err
			}
			if edirs != nil {
				emptys = append(emptys, edirs...)
			}
		}
	}
	return emptys, nil
}

// EmptyFloder 删除空目录
func EmptyFloder(dir string) error {
	emptys, err := FindEmptyFolder(dir)
	if err != nil {
		return err
	}
	for _, dir := range emptys {
		if err := os.Remove(dir); err != nil {
			return err
		} else {
			return err
		}
	}
	return nil
}

// substr 截取字符串的一部分
func substr(s string, pos, length int) (string, string) {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos : l+1]), string(runes[l+1:])
}

// GetParentDirectory 获取目录的父目录及余下路径
func GetParentDirectory(dirctory, prefix string) (string, string, bool) {
	index := strings.Index(dirctory, prefix)
	if index > -1 {
		res, eres := substr(dirctory, 0, index)
		return res, eres, true
	}
	return dirctory, "", false
}

// ReadLine 读取指定行的内容
func ReadLine(fileName string, lineNumber int) string {
	file, _ := os.Open(fileName)
	fileScanner := bufio.NewScanner(file)
	lineCount := 1
	for fileScanner.Scan() {
		if lineCount == lineNumber {
			return fileScanner.Text()
		}
		lineCount++
	}
	defer file.Close()
	return ""
}
