package helper

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qiaogw/pkg/s3cli"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/linxGnu/goseaweedfs"
	"github.com/qiaogw/pkg/config"
)

var DefaultConfig = &Config{}

type DirInfo struct {
	Path                  string
	Entries               []*Entries
	Limit                 int
	LastFileName          string
	ShouldDisplayLoadMore bool
}
type Entries struct {
	Mtime         time.Time   // time of last modification
	Crtime        time.Time   // time of creation (OS X only)
	Mode          os.FileMode // file mode
	Uid           uint32      // owner uid
	Gid           uint32      // group gid
	Mime          string      // mime type
	Replication   string      // replication
	Collection    string      // collection name
	TtlSec        int32       // ttl in seconds
	UserName      string
	GroupNames    []string
	SymlinkTarget string
	Md5           []byte
	FileSize      int64
	Extended      map[string][]byte
	FullPath      string
	Chunks        []*FileChunk `json:"chunks,omitempty"`
	IsDirectory   bool
	Type          string
	Name          string
	URLDir        string
	StorageClass string
	Restore string
}
type FileId struct {
	volume_id uint32
	file_key  uint64
	cookie    uint64
}
type FileChunk struct {
	FileId string `json:"file_id"`
	Offset int64  `json:"offset"`
	Size   int64  `json:"size"`
	Mtime  int64  `json:"mtime"`
	ETag   string `json:"e_tag"`
	Fid    FileId `json:"fid"`
}

type FilerURL struct {
	Public  string //Readonly URL
	Private string //Manage URL
}
type Config struct {
	Scheme    string
	Master    string
	Filers    []*FilerURL
	ChunkSize int64
	Timeout   time.Duration
	// TTL Time to live.
	// 3m: 3 minutes
	// 4h: 4 hours
	// 5d: 5 days
	// 6w: 6 weeks
	// 7M: 7 months
	// 8y: 8 years
	TTL string
}

func (c *Config) New() (*goseaweedfs.Seaweed, error) {
	if len(c.Scheme) == 0 {
		c.Scheme = "http"
	}
	if c.ChunkSize <= 0 {
		c.ChunkSize = 2 * 1024 * 1024
	}
	if c.Timeout <= 0 {
		c.Timeout = 5 * time.Minute
	}
	if len(c.Master) == 0 {
		c.Master = config.Config.Seaweed.Master
	}
	if c.Filers == nil || len(c.Filers) == 0 {
		c.Filers = []*FilerURL{
			{
				Public:  config.Config.Seaweed.Public,
				Private: config.Config.Seaweed.Private,
			},
		}
	}
	filers := make([]string, len(c.Filers))
	for index, filerURL := range c.Filers {
		filers[index] = filerURL.Private
	}
	return goseaweedfs.NewSeaweed(c.Master, filers, c.ChunkSize, &http.Client{Timeout: c.Timeout})
}

func (c *Config) MakeURL(path string, args url.Values) string {
	return MakeURL(c.Scheme, c.Master, path, args)
}

func MakeURL(scheme, host, path string, args url.Values) string {
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	if args != nil {
		u.RawQuery = args.Encode()
	}
	return u.String()
}

func (e *DirInfo) getSize() {
	for _, v := range e.Entries {
		v.GetSize()
	}
	return
}

func (v *Entries) GetSize() {
	//for _, v := range e.Entries {
	var size int64
	for _, m := range v.Chunks {
		size = size + m.Size
	}
	v.FileSize = size
	v.IsDirectory = v.Mode&os.ModeDir > 0
	v.Type = path.Ext(v.FullPath)
	if len(v.Type) > 0 {
		v.Type = v.Type[1:len(v.Type)]
	}
	v.Name = path.Base(v.FullPath)
	if v.IsDirectory {
		v.Type = "文件夹"
	}
	//}
	return
}

func GetInfo(pathStr, url string, v os.FileInfo) (e Entries, err error) {
	//var e helper.Entries
	e.Name = v.Name()
	e.FileSize = v.Size()
	e.IsDirectory = v.IsDir()
	e.Mode = v.Mode()
	e.Mtime = v.ModTime()
	e.FullPath = filepath.Join(pathStr, e.Name)
	e.URLDir = url + e.FullPath
	if e.IsDirectory {
		e.Type = "文件夹"
	} else {
		e.Type = path.Ext(e.Name)
		if len(e.Type)>0{
			e.Type = e.Type[1:len(e.Type)]
		}
	}
	sinfo:=v.Sys()
	rv := reflect.ValueOf(sinfo).Type()
	if rv.String()=="s3.Object"{
		ob:=sinfo.(s3.Object)
		e.StorageClass=*ob.StorageClass
		if *ob.StorageClass=="GLACIER" {
			svc:=s3cli.NewSvc()
			hb,err:=svc.Stat(svc.BucketName(),e.FullPath)
			if err==nil{
				if  hb.Restore==nil{
					e.Restore="未恢复"
				}else{
					if strings.Index(*hb.Restore, `ongoing-request="false"`)<0{
						e.Restore="恢复中"
					}else{
						e.Restore="已恢复"
					}
				}
			}
		}
	}
	// 取类型的元素
	//typeOfCat = typeOfCat.Elem()

	// 显示反射类型对象的名称和种类
	//logs.Infof("element name: '%v', element kind: '%v'\n", typeOfCat.Name(),typeOfCat.Kind())
	//field, b := reflect.TypeOf(sinfo).Elem().FieldByName(" StorageClass")
	//beego.Debug(field,b)
	if path.Ext(e.Name) == config.Config.LocalStore.Ignore {
		err = errors.New("隐藏")
	}
	return
}

func GetDirInfo(pathStr, url string, filelist []os.FileInfo) (e DirInfo) {
	//var e helper.Entries
	var infoList []*Entries
	var lastfile string
	for _, f := range filelist {
		e, err := GetInfo(pathStr, url, f)
		if err == nil {
			infoList = append(infoList, &e)
			lastfile = e.Name
		}
	}
	e.Path = pathStr
	e.Entries = infoList
	e.LastFileName = lastfile
	return
}
