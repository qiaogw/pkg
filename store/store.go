package store

import (
	"github.com/qiaogw/pkg/config"
	"github.com/qiaogw/pkg/store/driver"
)

var (
	// ErrExistsFile 文件不存在
	//ErrExistsFile = table.ErrExistsFile

	// BatchUpload 批量上传
	//BatchUpload = driver.BatchUpload

	// StorerRegister 存储引擎注册
	StorerRegister = driver.Register

	// StorerGet 获取存储引擎构造器
	StorerGet = driver.Get

	// StorerAll 存储引擎集合
	StorerAll = driver.All
	Store     driver.Storer
)

type (
	// Sizer 尺寸接口
	Sizer = driver.Sizer

	// Storer 文件存储引擎接口
	Storer = driver.Storer

	// Constructor 存储引擎构造函数
	Constructor = driver.Constructor
)

func Init() {
	Store, _ = StorerGet(config.Config.Store)
}
