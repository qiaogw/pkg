package local

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego"
	"github.com/qiaogw/pkg/config"
	"github.com/qiaogw/pkg/store/helper"
	"github.com/qiaogw/pkg/tools"

	//"os"
	"github.com/qiaogw/pkg/store"
)

const Name = `local`

var _ store.Storer = &Filesystem{}

func init() {
	store.StorerRegister(Name, func(typ string) (store.Storer, error) {
		beego.Debug(Name)
		return NewFilesystem(typ), nil
	})
}

func NewFilesystem(typ string, baseURLs ...string) *Filesystem {
	var host string
	addrs, err := net.InterfaceAddrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				host = ipnet.IP.String()
			}
		}
	}
	urlstr := fmt.Sprintf("http://%s:%v/%s", host, beego.BConfig.Listen.HTTPPort, "store")
	if len(baseURLs) > 0 {
		urlstr = baseURLs[0]
		urlstr = strings.TrimSuffix(urlstr, `/`)
	}
	beego.Debug(urlstr, beego.BConfig.Listen.HTTPPort, err)
	beego.Debug("s3 NewFilesystemNewFilesystemNewFilesystemNewFilesystem")
	return &Filesystem{
		Type:     typ,
		baseURL:  urlstr,
		BasePath: config.Config.LocalStore.Dir,
	}
}

// Filesystem 文件系统存储引擎
type Filesystem struct {
	Type     string
	baseURL  string
	BasePath string
}

func (f *Filesystem) Get(file string) (data io.Reader, err error) {
	panic("implement me")
}

func (f *Filesystem) Name() string {
	panic("implement me")
}

func (f *Filesystem) FileDir(subpath string) string {
	panic("implement me")
}

func (f *Filesystem) URLDir(subpath string) string {
	panic("implement me")
}

func (f *Filesystem) Put(dst string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	//创建本地文件
	dst = path.Join(f.BasePath, dst)
	if src == nil {
		dst = path.Join(dst, config.Config.LocalStore.Ignore)
	}
	err = tools.CheckPath(dst)
	if err != nil {
		return
	}
	destFile, err := os.Create(dst)
	if err != nil {
		return
	}
	//文件保存
	if src != nil {
		_, err = io.Copy(destFile, src)
	}
	if err != nil {
		return
	}
	//src.Close()
	destFile.Close()
	return
}

//func (f *Filesystem) Get(file string) (io.ReadCloser, error) {
//	panic("implement me")
//}

func (f *Filesystem) Exists(file string) (bool, error) {
	panic("implement me")
}

func (f *Filesystem) FileInfo(file string) (interface{}, error) {
	dir := filepath.Join(f.BasePath, file)
	info, err := os.Stat(dir)
	return info, err
}

func (f *Filesystem) PathInfo(pathStr string) (interface{}, error) {
	//die:=path.Dir(path1)
	dir := filepath.Join(f.BasePath, pathStr)
	fileInfoList, err := ioutil.ReadDir(dir)
	//var info helper.DirInfo
	//info.Path = path1
	//var filelist []*helper.Entries
	//for _, v := range fileInfoList {
	//	var e helper.Entries
	//	e.Name = v.Name()
	//	e.FileSize = v.Size()
	//	e.IsDirectory = v.IsDir()
	//	e.Mode = v.Mode()
	//	e.Mtime = v.ModTime()
	//	e.FullPath = filepath.Join(path1, e.Name)
	//	e.URLDir = f.baseURL + e.FullPath
	//	if e.IsDirectory {
	//		e.Type = "文件夹"
	//	} else {
	//		e.Type = path.Ext(e.Name)
	//		e.Type = e.Type[1:len(e.Type)]
	//	}
	//	if path.Ext(e.Name) != config.Config.LocalStore.Ignore {
	//		filelist = append(filelist, &e)
	//	}
	//}
	//if err != nil {
	//	beego.Error(err)
	//}
	//info.Entries = filelist
	info := helper.GetDirInfo(pathStr, f.BaseURL(), fileInfoList)
	dTree := make(map[string]interface{})
	dTree["info"] = info
	lt := tools.GetDirList(filepath.Join(f.BasePath, "/"), "/")
	var nt tools.DirBody
	nt.Dir = "/"
	nt.Label = "全部文件夹"
	nt.Children = lt
	nt.Icon = "el-icon-folder-add"
	dTree["tree"] = nt
	return dTree, err
	//panic("implement me")
}

func (s *Filesystem) GetDir(dir string) (interface{}, error) {
	panic("implement me")
}
func (f *Filesystem) SendFile(file string) error {
	panic("implement me")
}

func (f *Filesystem) Delete(file string) error {
	dst := path.Join(f.BasePath, file)
	err := os.Remove(dst)
	return err
}

func (f *Filesystem) DeleteDir(dir string) error {
	dst := path.Join(f.BasePath, dir)
	err := os.RemoveAll(dst)
	return err
}

func (f *Filesystem) Move(src, dst string) error {
	srcBase := path.Join(f.BasePath, src)
	dstBase := path.Join(f.BasePath, dst)
	err := os.Rename(srcBase, dstBase)
	return err
}

func (f *Filesystem) PublicURL(dst string) string {
	panic("implement me")
}

func (f *Filesystem) URLToFile(viewURL string) string {
	panic("implement me")
}

func (f *Filesystem) URLToPath(viewURL string) string {
	panic("implement me")
}

func (f *Filesystem) SetBaseURL(baseURL string) {
	panic("implement me")
}

func (f *Filesystem) BaseURL() string {
	return f.baseURL
}

func (f *Filesystem) FixURL(content string, embedded ...bool) string {
	panic("implement me")
}

func (f *Filesystem) FixURLWithParams(content string, values url.Values, embedded ...bool) string {
	panic("implement me")
}

func (f *Filesystem) Close() error {
	panic("implement me")
}

func (f *Filesystem) Search(seachStr string) interface{} {
	//panic("implement me")
	filelist := make([]*helper.Entries, 0)
	fs := helper.SearchFileList(filepath.Join(f.BasePath, "/"), "/", seachStr, f.baseURL, config.Config.LocalStore.Ignore, filelist)
	dTree := make(map[string]interface{})
	var info helper.DirInfo
	info.Path = "/搜索:" + seachStr
	info.Entries = fs
	dTree["info"] = info
	lt := tools.GetDirList(filepath.Join(f.BasePath, "/"), "/")
	var nt tools.DirBody
	nt.Dir = "/"
	nt.Label = "全部文件夹"
	nt.Children = lt
	nt.Icon = "el-icon-folder-add"
	dTree["tree"] = nt
	return dTree
	//return lt
}
