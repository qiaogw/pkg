package redis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	RedisPool *redis.Pool
	//DefautConfig RedisConfig
)

type RedisConfig struct {
	Enable      bool          `label:"是否启用"`
	Key         string        `label:"Redis collection 的名称"`
	Addr        string        `label:"地址"`
	Password    string        `label:"密码"`
	DBNum       int           `label:"数据库"`
	MaxActive   int           `label:"最大连接数"`
	MaxIdle     int           `label:"最大空闲连接数"`
	IdleTimeout time.Duration `label:"空闲连接超时时间"`
	Wait        bool          `label:"如果超过最大连接，是报错，还是等待。"`
}

func Init(conf RedisConfig) {
	if conf.IdleTimeout < 1000 {
		conf.IdleTimeout = 1000
	}
	RedisPool = &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: conf.IdleTimeout,
		Wait:        conf.Wait,
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			setdb := redis.DialDatabase(conf.DBNum)
			setPasswd := redis.DialPassword(conf.Password)
			setWriteTimeOut := redis.DialWriteTimeout(conf.IdleTimeout)
			setReadTimeOut := redis.DialReadTimeout(conf.IdleTimeout)
			SetConnetTimeOut := redis.DialConnectTimeout(conf.IdleTimeout)
			c, err := redis.Dial("tcp", conf.Addr, setPasswd, setdb, setReadTimeOut, setWriteTimeOut, SetConnetTimeOut)
			if err != nil {
				return nil, err
			}
			//if _, err := c.Do("AUTH", conf.Password); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			if _, err := c.Do("SELECT", conf.DBNum); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func ListAdd(vName string, evt interface{}) error {
	c := RedisPool.Get()
	defer c.Close()
	_, err := c.Do("RPush", vName, evt)
	return err
}

func ListRemove(vName string, evt interface{}) error {
	c := RedisPool.Get()
	defer c.Close()
	_, err := c.Do("LRem", vName, 0, evt)
	return err
}

func ListGet(vName string, begin, end int, ul interface{}) (err error) {
	c := RedisPool.Get()
	defer c.Close()
	values, err := redis.Values(c.Do("LRange", vName, begin, end))
	err = redis.ScanSlice(values, ul)
	return
}
func ListStructAdd(vName string, evt interface{}) error {
	c := RedisPool.Get()
	defer c.Close()
	se, _ := json.Marshal(evt)
	_, err := c.Do("RPush", vName, se)
	return err
}

func ListStructGet(vName string, begin, end int) (interface{}, error) {
	c := RedisPool.Get()
	defer c.Close()
	values, err := redis.Strings(c.Do("LRange", vName, begin, end))
	var et interface{}
	var rl []interface{}
	for _, v := range values {
		err = json.Unmarshal([]byte(v), &et)
		rl = append(rl, et)
	}
	return rl, err
}
func ListStructRemove(vName string, evt interface{}) error {
	se, _ := json.Marshal(evt)
	return ListRemove(vName, se)

}
func StructAdd(vName string, evt interface{}) error {
	c := RedisPool.Get()
	defer c.Close()
	_, err := c.Do("HMSET", redis.Args{}.Add(vName).AddFlat(evt)...)
	//c.Flush()
	//_, err := c.Receive()
	return err
}

func StructGet(vName string, evt interface{}) error {
	c := RedisPool.Get()
	defer c.Close()
	_, err := c.Do("HMSET", redis.Args{}.Add(vName).AddFlat(evt)...)
	//c.Flush()
	//_, err := c.Receive()
	return err
}
func Remove(vName string) error {
	c := RedisPool.Get()
	defer c.Close()
	_, err := c.Do("del", vName)
	//c.Flush()
	//_, err := c.Receive()
	return err
}

func StructGetAll(vName interface{}, ul interface{}) (err error) {
	c := RedisPool.Get()
	defer c.Close()
	result, err := redis.Values(c.Do("hgetall", vName))
	err = redis.ScanStruct(result, ul)
	//ul, err = redis.StringMap(c.Do("hgetall", vName))
	return
}

//func LGet(list string, id interface{}) (evt interface{}, err error) {
//	c := RedisPool.Get()
//	defer c.Close()
//	c.Send("HMSET", id, evt)
//	c.Send("RPush", list, id)
//	c.Flush()
//	v, err := c.Receive()
//	beego.Debug(v, err)
//	return err
//}
