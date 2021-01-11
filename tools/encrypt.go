package tools

import (
	"crypto/md5"
	"math/rand"

	//	"crypto/rand"
	//	"encoding/base64"
	"encoding/hex"
	//	"io"

	"reflect"
	//	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	//	beego "github.com/beego/beego/v2/adapter"
)

//Md5 把字符串转换为md5方法
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5 把字节转化为MD5字符串
func MD5(b []byte) string {
	vCrypto := md5.New()
	vCrypto.Write(b)
	return hex.EncodeToString(vCrypto.Sum(nil))
}

// EncodeUserPwd 加密用户密码
func EncodeUserPwd(uname, pwd, salt string) string {
	return MD5([]byte(strings.Join([]string{uname, "$user$", pwd}, salt)))
}

// EncodeMemberPwd 加密客户密码
func EncodeMemberPwd(uname, pwd string) string {
	return MD5([]byte(strings.Join([]string{uname, "$member$", pwd}, "")))
}

// EncodeToken 加密token
func EncodeToken(sessionid string) string {
	timez := time.Now().String()
	return MD5([]byte(strings.Join([]string{sessionid, "$token$", timez}, "")))
}

// Struct2Map stuct 转 map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//GetGuid 获取Guid方法，获取48位uuid
func GetGuid() string {
	uid := uuid.NewV1()
	return uid.String()
}

// 生成长度为length的随机字符串
func RandString(length int64) string {
	sources := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sourceLength := len(sources)
	var i int64 = 0
	for ; i < length; i++ {
		result = append(result, sources[r.Intn(sourceLength)])
	}

	return string(result)
}
