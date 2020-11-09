package seaweed

import (
	"bytes"
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/qiaogw/pkg/config"

	//"io/ioutil"
	"strings"

	//"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/linxGnu/goseaweedfs"
	"github.com/qiaogw/pkg/store"
	"github.com/qiaogw/pkg/store/driver/local"
	"github.com/qiaogw/pkg/store/helper"
)

const Name = `seaweedfs`

var _ store.Storer = &Seaweedfs{}

func init() {
	store.StorerRegister(Name, func(typ string) (store.Storer, error) {
		return NewSeaweedfs(typ), nil
	})
}

func NewSeaweedfs(typ string) *Seaweedfs {
	a, _ := helper.DefaultConfig.New()
	return &Seaweedfs{
		config:     helper.DefaultConfig,
		instance:   a,
		Filesystem: local.NewFilesystem(typ),
		baseURL:    config.Config.Seaweed.Public,
	}
}

type Seaweedfs struct {
	config   *helper.Config
	instance *goseaweedfs.Seaweed
	*local.Filesystem
	baseURL string
}

func (s *Seaweedfs) FileDir(subpath string) string {
	return path.Join(helper.UploadURLPath, s.Type, subpath)
}

func (s *Seaweedfs) URLDir(subpath string) string {
	return path.Join(helper.UploadURLPath, s.Type, subpath)
}

func (s *Seaweedfs) Put(dst string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	filer := s.instance.Filers()
	res, err := filer[0].Upload(src, size, dst, "col", s.config.TTL)
	return res.FileID, res.FileURL, err
}

func (s *Seaweedfs) Get(file string) (data io.Reader, err error) {
	//filer := s.instance.Filers()[0]
	//_, _, err = filer.Get(file, nil, nil)
	//err = filer.Download(file, nil, func(r io.Reader) error {
	//	data = r
	//	return err
	//})
	return data, err
}

func (s *Seaweedfs) Exists(file string) (bool, error) {
	filer := s.instance.Filers()[0]
	data, status, err := filer.Get(file, nil, nil)
	if status == 200 && len(data) > 0 {
		return true, nil
	}
	return false, err
}

func (s *Seaweedfs) FilerInfo(file string) (os.FileInfo, error) {
	filer := s.instance.Filers()[0]
	args := make(url.Values)
	args.Set("pretty", "y")
	header := make(map[string]string)
	header["Accept"] = "application/json"
	data, _, err := filer.Get(file, nil, nil)
	beego.Debug(err, s.BaseURL(), s.baseURL)
	var fileInfo os.FileInfo
	//fileInfo := make(map[string]interface{})
	err = json.Unmarshal(data, &fileInfo)
	return fileInfo, err
}

func (s *Seaweedfs) FileInfo(file string) (interface{}, error) {
	var fileList helper.DirInfo
	var fileInfo helper.Entries
	dir := filepath.Dir(file)
	filer := s.instance.Filers()[0]
	header := make(map[string]string)
	header["Accept"] = "application/json"
	args := make(url.Values)
	args.Set("pretty", "y")
	data, _, err := filer.Get(dir, args, header)
	if err != nil {
		return fileInfo, err
	}
	err = json.Unmarshal(data, &fileList)
	if err != nil {
		return fileInfo, err
	}
	for _, v := range fileList.Entries {
		if file == v.FullPath {
			v.GetSize()
			fileInfo = *v
			break
		}
	}
	return fileInfo, err
}

//type dirinfo struct {
//	driver.DirInfo
//}

func (s *Seaweedfs) PathInfo(path string) (interface{}, error) {
	var info helper.DirInfo
	//dir := filepath.Dir(file)
	filer := s.instance.Filers()[0]
	header := make(map[string]string)
	header["Accept"] = "application/json"
	args := make(url.Values)
	args.Set("pretty", "y")
	data, _, err := filer.Get(path, args, header)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(data, &info)
	//info.GetSize()
	return info, err
}

func (s *Seaweedfs) GetDir(dir string) (interface{}, error) {
	panic("implement me")
}

func (s *Seaweedfs) SendFile(file string) error {
	panic("implement me")
}

func (s *Seaweedfs) Delete(file string) error {
	filer := s.instance.Filers()[0]
	err := filer.Delete(file, nil)
	return err
}

func (s *Seaweedfs) DeleteDir(dir string) error {
	filer := s.instance.Filers()[0]
	args := make(url.Values)
	args.Set("recursive", "true")
	args.Set("ignoreRecursiveError", "true")
	err := filer.Delete(dir, args)
	return err
}

func (s *Seaweedfs) Move(src, dst string) error {
	return s.Rename(src, dst)
}

func (s *Seaweedfs) Rename(src, dst string) error {
	//panic("implement me")
	filer := s.instance.Filers()[0]
	data, _, err := filer.Get(src, nil, nil)
	if err != nil {
		return err
	}
	//var fi io.Reader
	fi := bytes.NewReader(data)
	//binary.Read(bytes.NewReader(data), binary.LittleEndian, &fi)
	_, err = filer.Upload(fi, int64(len(data)), dst, "col", s.config.TTL)
	if err != nil {
		return err
	}
	err = filer.Delete(src, nil)
	return err
}

func (s *Seaweedfs) URLToFile(viewURL string) string {
	dstFile := s.URLToPath(viewURL)
	dstFile = strings.TrimPrefix(dstFile, strings.TrimRight(s.URLDir(``), `/`)+`/`)
	return dstFile
}

func (s *Seaweedfs) URLToPath(viewURL string) string {
	if len(s.baseURL) > 0 {
		viewURL = strings.TrimPrefix(viewURL, s.baseURL+`/`)
		if !strings.HasPrefix(viewURL, `/`) {
			viewURL = `/` + viewURL
		}
	}
	return viewURL
}

func (s *Seaweedfs) SetBaseURL(baseURL string) {
	baseURL = strings.TrimSuffix(baseURL, `/`)
	s.baseURL = baseURL
}

func (s *Seaweedfs) BaseURL() string {
	return s.baseURL
}

func (s *Seaweedfs) Name() string {
	return Name
}

func (s *Seaweedfs) filepath(fname string) string {
	return s.URLDir(fname)
}

func (s *Seaweedfs) PublicURL(dstFile string) string {
	return s.config.Filers[0].Public + dstFile
}

func (s *Seaweedfs) FixURL(content string, embedded ...bool) string {
	if len(embedded) > 0 && embedded[0] {
		return helper.ReplaceAnyFileName(content, func(r string) string {
			return s.PublicURL(r)
		})
	}
	return s.PublicURL(content)
}

func (s *Seaweedfs) FixURLWithParams(content string, values url.Values, embedded ...bool) string {
	if len(embedded) > 0 && embedded[0] {
		return helper.ReplaceAnyFileName(content, func(r string) string {
			return s.FixURLWithParams(s.PublicURL(r), values)
		})
	}
	return s.FixURLWithParams(s.PublicURL(content), values)
}
func (f *Seaweedfs) Search(seachStr string) interface{} {
	panic("implement me")
}
