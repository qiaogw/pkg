package redis

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"github.com/qiaogw/pkg/config"
)

var RedisPool *redis.Pool

func init() {
	//pool = &redis.Pool{
	//	MaxIdle:     16,
	//	MaxActive:   1024,
	//	IdleTimeout: 300,
	//	Dial: func() (redis.Conn, error) {
	//		return redis.Dial("tcp", "localhost:6379")
	//	},
	//}
	RedisPool = &redis.Pool{
		MaxIdle:     config.Config.Redis.MaxIdle,
		MaxActive:   config.Config.Redis.MaxActive,
		IdleTimeout: config.Config.Redis.IdleTimeout,
		Wait:        config.Config.Redis.Wait,
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Config.Redis.Addr)
			if err != nil {
				return nil, err
			}
			//if _, err := c.Do("AUTH", config.Config.Redis.Password); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			if _, err := c.Do("SELECT", config.Config.Redis.DBNum); err != nil {
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
	beego.Debug(vName, evt)
	_, err := c.Do("RPush", vName, evt)
	//c.Flush()
	//_, err := c.Receive()
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
	//id := evt.(map[string]int)["Id"]
	//ids := strconv.Itoa(id)
	//_, err := c.Do("HMSET", redis.Args{}.Add(vName+"_"+ids).AddFlat(evt)...)
	se, _ := json.Marshal(evt)
	_, err := c.Do("RPush", vName, se)
	//c.Flush()
	//_, err := c.Receive()
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
	beego.Debug(vName)
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
