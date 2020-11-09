// Copyright 2018 cloudy itcloudy@qq.com.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// IntToStr converts integer to string
func IntToStr(num int) string {
	return strconv.Itoa(num)
}

// StrToInt64 converts string to int64
func StrToInt64(s string) int64 {
	int64, _ := strconv.ParseInt(s, 10, 64)
	return int64
}

// BytesToInt64 converts []bytes to int64
func BytesToInt64(s []byte) int64 {
	int64, _ := strconv.ParseInt(string(s), 10, 64)
	return int64
}

// StrToUint64 converts string to the unsinged int64
func StrToUint64(s string) uint64 {
	ret, _ := strconv.ParseUint(s, 10, 64)
	return ret
}

// Float64ToStr converts float64 to string
func Float64ToStr(f float64, prec ...string) string {
	return strconv.FormatFloat(f, 'f', 13, 64)
}

// StrToInt converts string to integer
func StrToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
func StringToIntDefault(str string, defVal int) int {
	if in, err := strconv.Atoi(str); err != nil {
		return defVal
	} else {
		return in
	}
}

// StrToFloat64 converts string to float64
func StrToFloat64(s string) float64 {
	Float64, _ := strconv.ParseFloat(s, 64)
	return Float64
}

// BytesToFloat64 converts []byte to float64
func BytesToFloat64(s []byte) float64 {
	Float64, _ := strconv.ParseFloat(string(s), 64)
	return Float64
}

// BytesToInt converts []byte to integer
func BytesToInt(s []byte) int {
	i, _ := strconv.Atoi(string(s))
	return i
}
func StrToBool(s string) bool {
	if s == "0" || s == "false" {
		return false
	}
	return true
}

//SnakeString 蛇形转换 string, XxYy to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

//CamelString 驼峰转换 string ， xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// camelCase 驼峰转换 string，converts a _ delimited string to camel case
// e.g. very_important_person => VeryImportantPerson
func CamelCase(in string) string {
	tokens := strings.Split(in, "_")
	for i := range tokens {
		tokens[i] = strings.Title(strings.Trim(tokens[i], " "))
	}
	return strings.Join(tokens, "")
}

// formatSourceCode formats source files
func FormatSourceCode(filename string) {
	cmd := exec.Command("gofmt", "-w", filename)
	cmd.Run()
}

// TypeInt 断言获取int或int64为int,其他类型为0
func TypeInt(num interface{}) int {
	n32, ok := num.(int)
	if !ok {
		n64, ok := num.(int64)
		if ok {
			n32 = int(n64)
		} else {
			n32 = 0
		}
	}
	return n32
}

// TypeInt64 断言获取int或int64为int64,其他类型为0
func TypeInt64(num interface{}) int64 {
	n64, ok := num.(int64)
	if !ok {
		n32, ok := num.(int)
		if ok {
			n64 = int64(n32)
		} else {
			n64 = int64(0)
		}
	}
	return n64
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// GetConfigFromPath read config from path and returns  struct
func GetConfigFromPath(path string, vs interface{}) (err error) {
	//log.WithFields(log.Fields{"path": path}).Info("Loading config")

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return errors.Errorf("Unable to load config file %s", path)
	}

	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "reading config")
	}

	//c := new(interface{})
	err = viper.Unmarshal(vs)
	if err != nil {
		err = errors.Wrapf(err, "marshalling config to global struct variable")
	}

	return
}
