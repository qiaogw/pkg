package driver

import (
	"io"

	"net/url"
	"sort"

	"github.com/fatih/color"
	"github.com/qiaogw/pkg/logs"

	"github.com/qiaogw/pkg/errors"
)

var (
	// ErrExistsFile 文件不存在
	ErrExistsFile = errors.ErrExistsFile
)

type Result struct {
	FileID   int64
	FileName string
	FileURL  string
	FileType string
	FileSize int64
	SavePath string
	Md5      string
	Addon    interface{}
	// fileNameGenerator func(string) (string, error)
}

// Sizer 尺寸接口
type Sizer interface {
	Size() int64
}

// Storer 文件存储引擎接口
type Storer interface {
	// 引擎名
	Name() string

	// FileDir 文件夹物理路径
	FileDir(subpath string) string

	// URLDir 文件夹网址路径
	URLDir(subpath string) string

	// Put 保存文件
	Put(dst string, src io.Reader, size int64) (savePath string, viewURL string, err error)

	// Get 获取文件
	Get(file string) (data io.Reader, err error)

	// Exists 文件是否存在
	Exists(file string) (bool, error)

	// FileInfo 文件信息
	FileInfo(file string) (interface{}, error)

	// PathInfo 文件夹信息
	PathInfo(path string) (interface{}, error)

	//GetDir 获取文件夹树
	GetDir(dir string) (interface{}, error)

	// Delete 删除文件
	Delete(file string) error

	// DeleteDir 删除目录
	DeleteDir(dir string) error

	// Move 移动文件
	Move(src, dst string) error

	// PublicURL 文件物理路径转网址
	PublicURL(dst string) string

	// URLToFile 网址转文件存储路径(非完整路径)
	URLToFile(viewURL string) string

	// URLToPath 网址转文件路径(完整路径)
	URLToPath(viewURL string) string

	// 根网址(末尾不含"/")
	SetBaseURL(baseURL string)
	Search(seachStr string) interface{}
	BaseURL() string

	// FixURL 修正网址
	FixURL(content string, embedded ...bool) string

	// FixURLWithParams 修正网址并增加网址参数
	FixURLWithParams(content string, values url.Values, embedded ...bool) string
}

// Constructor 存储引擎构造函数
type Constructor func(typ string) (Storer, error)

var storers = map[string]Constructor{}

// DefaultConstructor 默认构造器
var DefaultConstructor Constructor

// Register 存储引擎注册
func Register(engine string, constructor Constructor) {
	logs.Info(color.CyanString(`storer.register:`), engine)
	storers[engine] = constructor
}

// Get 获取存储引擎构造器
func Get(engine string) (Storer, error) {
	constructor, ok := storers[engine]
	if !ok {
		return DefaultConstructor(engine)
	}
	return constructor(engine)
}

//// GetBySettings 获取存储引擎构造器
//func GetBySettings() Constructor {
//	engine := `local`
//	//storerConfig, ok := storer.GetOk()
//	//if ok {
//	//	engine = storerConfig.Name
//	//}
//	return Get(engine)
//}

// All 存储引擎集合
func All() map[string]Constructor {
	return storers
}

// AllNames 存储引擎集合
func AllNames() []string {
	var names []string
	for name := range storers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
