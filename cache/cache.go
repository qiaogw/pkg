//Package cache 功能
package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/astaxie/beego/cache/ssdb"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
	"github.com/qiaogw/pkg/config"
	"time"
)

var cc cache.Cache

// InitCache 初始化cache
func InitCache() {
	var err error
	cacheConfig := config.Config.Cache.CacheType
	cc = nil
	if len(cacheConfig) < 4 {
		cacheConfig = "memory"
	}
	beego.Notice("Cache 使用:", cacheConfig)
	switch cacheConfig {
	case "ssdb":
		err = initSsdb()
	case "redis":
		err = initRedis()
	case "memcache":
		err = initMemcache()
	case "file":
		err = initFile()
	case "memory":
		err = initMemory()
	default:
		err = errors.Errorf("Cache driver is not allowed:", cacheConfig)
		beego.Error(err)
	}
	if err != nil {
		err = initMemory()
	}
}

func initFile() (err error) {
	cc, err = cache.NewCache("File", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"120"}`)

	if err != nil {
		beego.Error(err)
	}
	return
}

func initMemory() (err error) {
	cc, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		beego.Error(err)
	}
	return
}
func initMemcache() (err error) {
	cc, err = cache.NewCache("memcache", `{"conn":"`+beego.AppConfig.String("memcache_host")+`"}`)

	if err != nil {
		beego.Error(err)
	}
	return
}

func initRedis() (err error) {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("initcacheredis err is :", r)
			cc = nil
		}
	}()
	key := "cacheCollectionName" //config.Config.Redis.Key
	conn := config.Config.Redis.Addr
	dbNum := 0 //config.Config.Redis.DBNum
	password := config.Config.Redis.Password
	confStr := fmt.Sprintf(`{"key":"%v","conn":"%v","dbNum":"%v","password":"%v"}`, key, conn, dbNum, password)
	cc, err = cache.NewCache("redis", confStr)
	if err != nil {
		beego.Error("redis连接失败！失败原因：", confStr, err)
	} else {
		beego.Notice("redis连接成功: ", confStr)
	}
	return
}
func initSsdb() (err error) {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("initcachessdb err is :", r)
			//fmt.Println("initial redis error caught: %v\n", r)
			cc = nil
		}
	}()
	cc, err = cache.NewCache("ssdb", `{"conn":"`+beego.AppConfig.String("ssdb_host")+":"+
		beego.AppConfig.String("ssdb_port")+`"}`)
	if err != nil {
		beego.Error("ssdb连接失败！失败原因：", err)
	} else {
		beego.Notice("ssdb连接成功: %s:%s", beego.AppConfig.String("ssdb_host"), beego.AppConfig.String("ssdb_port"))
	}
	return
}

// SetCache 插入cache
func SetCache(key string, value interface{}, ts ...int64) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("set cache error caught: %v\n", r)
			cc = nil
		}
	}()

	timeout := config.Config.Cache.CacheExpire
	timeouts := time.Duration(timeout) * time.Second
	if len(ts) > 0 {
		timeouts = time.Duration(ts[0]) * time.Second
	}
	// beego.Error(data)

	err = cc.Put(key, data, timeouts)
	if err != nil {
		//fmt.Println("Cache失败，key:", key)
		return err
	}
	return nil

}

//GetCache 获取cache
func GetCache(key string, to interface{}) error {
	if cc == nil {
		// return errors.New("cc is nil")
		InitCache()
	}
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()

	data := cc.Get(key)
	// beego.Error(data)
	if data == nil {
		return errors.New("Cache不存在")
	}
	// log.Pinkln(data)

	err := Decode(data.([]byte), to)
	if err != nil {
		//fmt.Println("获取Cache失败", key, err)
	}

	return err
}

// DelCache 产出cache
func DelCache(key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()

	err := cc.Delete(key)
	if err != nil {
		return errors.New("Cache删除失败")
	}
	return nil

}

// Encode 用gob进行数据编码
func Encode(data interface{}) ([]byte, error) {
	gob.Register(time.Time{})
	gob.Register(orm.ParamsList{})
	gob.Register(map[string]interface{}{})
	gob.Register(data)
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode 用gob进行数据解码
func Decode(data []byte, to interface{}) error {
	gob.Register(time.Time{})
	//gob.Register(orm.ParamsList{})
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
