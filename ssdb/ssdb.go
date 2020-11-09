package ssdb

//ssdb连接池
import (
	//	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/ssdb/gossdb/ssdb"
)

//ssdb连接信息
type SsdbProvider struct {
	conn *ssdb.Client
	Host string
	Port int
	//	MaxLifetime int64
}

//ssdb初始化连接
func (p *SsdbProvider) connectInit() error {
	var err error
	p.Host = beego.AppConfig.String("ssdb_host")
	port := beego.AppConfig.String("ssdb_port")
	p.Port, _ = strconv.Atoi(port)
	p.conn, err = ssdb.Connect(p.Host, p.Port)

	return err
}

////NewSsdbCache create new ssdb adapter.
//func NewSsdbCache() cache.Cache {
//	return &Cache{}
//}

// Get get value from memcache.
func (rc *SsdbProvider) Get(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
		}
	}
	value, err := rc.conn.Get(key)
	if err == nil {
		return value
	}
	return nil
}

// GetMulti get value from memcache.
func (rc *SsdbProvider) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var values []interface{}
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			for i := 0; i < size; i++ {
				values = append(values, err)
			}
			return values
		}
	}
	res, err := rc.conn.Do("multi_get", keys)
	resSize := len(res)
	if err == nil {
		for i := 1; i < resSize; i += 2 {
			values = append(values, res[i+1])
		}
		return values
	}
	for i := 0; i < size; i++ {
		values = append(values, err)
	}
	return values
}

// DelMulti get value from memcache.
func (rc *SsdbProvider) DelMulti(keys []string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Do("multi_del", keys)
	return err
}

// Put put value to memcache. only support string.
func (rc *SsdbProvider) Put(key string, value interface{}, timeout time.Duration) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	v, ok := value.(string)
	if !ok {
		return errors.New("value must string")
	}
	var resp []string
	var err error
	ttl := int(timeout / time.Second)
	if ttl < 0 {
		resp, err = rc.conn.Do("set", key, v)
	} else {
		resp, err = rc.conn.Do("setx", key, v, ttl)
	}
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return errors.New("bad response")
}

// Delete delete value in memcache.
func (rc *SsdbProvider) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Del(key)
	return err
}

// Incr increase counter.
func (rc *SsdbProvider) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Do("incr", key, 1)
	return err
}

// Decr decrease counter.
func (rc *SsdbProvider) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Do("incr", key, -1)
	return err
}

// IsExist check value exists in memcache.
func (rc *SsdbProvider) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	resp, err := rc.conn.Do("exists", key)
	if err != nil {
		return false
	}
	if len(resp) == 2 && resp[1] == "1" {
		return true
	}
	return false

}

// ClearAll clear all cached in memcache.
func (rc *SsdbProvider) ClearAll() error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	keyStart, keyEnd, limit := "", "", 50
	resp, err := rc.Scan(keyStart, keyEnd, limit)
	for err == nil {
		size := len(resp)
		if size == 1 {
			return nil
		}
		keys := []string{}
		for i := 1; i < size; i += 2 {
			keys = append(keys, resp[i])
		}
		_, e := rc.conn.Do("multi_del", keys)
		if e != nil {
			return e
		}
		keyStart = resp[size-2]
		resp, err = rc.Scan(keyStart, keyEnd, limit)
	}
	return err
}

// Scan key all cached in ssdb.
func (rc *SsdbProvider) Scan(keyStart string, keyEnd string, limit int) ([]string, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil, err
		}
	}
	resp, err := rc.conn.Do("scan", keyStart, keyEnd, limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
