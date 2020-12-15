package s3

import (
	//"fmt"
	//"github.com/astaxie/beego"

	"github.com/aws/aws-sdk-go/service/s3"
	"time"

	"github.com/pkg/errors"
	"github.com/qiaogw/pkg/s3cli"
	"github.com/qiaogw/pkg/store"
	"github.com/qiaogw/pkg/store/helper"
	"github.com/qiaogw/pkg/tools"
	//"github.com/astaxie/beego/context"
	"io"
	"net/url"
	"os"
	"path"
)

const Name = `s3`

var _ store.Storer = &Filesystem{}

func init() {
	store.StorerRegister(Name, func(typ string) (store.Storer, error) {
		return NewFilesystem(typ), nil
	})
}

func NewFilesystem(typ string, baseURLs ...string) *Filesystem {
	//var host string
	//addrs, _ := net.InterfaceAddrs()
	//for _, address := range addrs {
	//	// 检查ip地址判断是否回环地址
	//	if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	//		if ipnet.IP.To4() != nil {
	//			host = ipnet.IP.String()
	//		}
	//	}
	//}
	//urlstr := fmt.Sprintf("http://%s:%v/%s", host, beego.BConfig.Listen.HTTPPort, "store")
	//beego.Debug("s3 NewFilesystemNewFilesystemNewFilesystemNewFilesystem", urlstr)
	mgr, _ := New(1024)
	svc := s3cli.NewSvc()
	return &Filesystem{
		Type: typ,
		BasePath: "/",
		mgr:      mgr,
		svc:      svc,
	}
}

// Filesystem 文件系统存储引擎
type Filesystem struct {
	Type     string
	baseURL  string
	BasePath string
	mgr      *S3Manager
	svc      *s3cli.Svc
}

func (f *Filesystem) Get(file string) (data io.Reader, err error) {
	//data, err = f.mgr.Get(file)
	//if err != nil {
	//	err = errors.WithMessage(err, Name)
	//}
	data,err=f.svc.Get(f.svc.BucketName(),file)
	return
}

func (f *Filesystem) Name() string {
	return Name
}

func (f *Filesystem) FileDir(subpath string) string {
	return f.BasePath
}

func (f *Filesystem) URLDir(subpath string) string {
	return f.baseURL
}

func (f *Filesystem) Put(dst string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	err=f.svc.Put(src,dst)
	return
}


func (f *Filesystem) Exists(fileStr string) (bool, error) {
	return f.svc.Exists(fileStr)
}

func (f *Filesystem) FileInfo(fileStr string) (interface{}, error) {
	objectInfo, err := f.svc.Stat(f.svc.BucketName(),fileStr)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	var finfo s3.Object
	*finfo.Key=fileStr
	finfo.ETag=objectInfo.ETag
	finfo.LastModified=objectInfo.LastModified
	finfo.StorageClass=objectInfo.StorageClass
	return s3cli.NewFileInfo(finfo), nil
}

func (f *Filesystem) PathInfo(pathStr string) (interface{}, error) {
	err, _, objectInfo := f.svc.List(f.mgr.bucketName, pathStr,true)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	pathinfo := helper.GetDirInfo(pathStr, f.BaseURL(), objectInfo)
	fTree := f.mgr.ListTree(f.BasePath)
	var nt tools.DirBody
	nt.Dir = "/"
	nt.Label = "全部文件夹"
	nt.Children = fTree
	nt.Icon = "el-icon-folder-add"
	dTree := make(map[string]interface{})
	dTree["tree"] = nt
	dTree["info"] = pathinfo
	return dTree, nil
}

func (f *Filesystem) GetDir(dir string) (interface{}, error) {
	var dirs []os.FileInfo
	err:=f.svc.GetDir(&dirs,f.svc.BucketName(),dir,"",time.Now())
	return dirs,err
}
// SendFile 下载文件
func (f *Filesystem) SendFile(file string) error {
	_, err := f.Get(file)
	if err != nil {
		return errors.WithMessage(err, Name)
	}
	return err
}

func (f *Filesystem) Delete(file string) error {
	err:=f.svc.RemoveObject(f.svc.BucketName(),file)
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
	return nil
}

func (f *Filesystem) Search(seachStr string) interface{} {
	err, _, objectInfo := f.svc.List(f.mgr.bucketName, seachStr,false)
	if err != nil {
		return nil
	}
	pathinfo := helper.GetDirInfo(seachStr, f.BaseURL(), objectInfo)
	fTree := f.mgr.ListTree(f.BasePath)
	var nt tools.DirBody
	nt.Dir = "/"
	nt.Label = "全部文件夹"
	nt.Children = fTree
	nt.Icon = "el-icon-folder-add"
	dTree := make(map[string]interface{})
	dTree["tree"] = nt
	dTree["info"] = pathinfo
	return dTree
}
