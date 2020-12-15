package tools

import (
	"bytes"
	"crypto/rand"
	"math/big"

	uuid "github.com/satori/go.uuid"

	"strconv"
	"strings"
	"time"
)

//Substr 截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//Substr2 截取字符串 start 起点下标 end 终点下标(不包括)
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 {
		start = 0
	}
	if start > length {
		start = length
	}
	if end < 0 {
		end = 0
	}
	if end > length {
		end = length
	}

	return string(rs[start:end])
}

//GetBetweenStr 获取字符串 str在substr之后的
func GetBetweenStr(str, substr string) string {
	n := strings.Index(str, substr)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	return str
}

// RandInt64 取随64位机数
func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

// Strim 去除string空
func Strim(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

// Unicode 取字符串unicode
func Unicode(rs string) string {
	json := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			json += string(r)
		} else {
			json += "\\u" + strconv.FormatInt(int64(rint), 16)
		}
	}
	return json
}

// HTMLEncode html格式编码
func HTMLEncode(rs string) string {
	html := ""
	for _, r := range rs {
		html += "&#" + strconv.Itoa(int(r)) + ";"
	}
	return html
}

//RemoveDuplicatesAndEmpty string数组去重
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	aLen := len(a)
	for i := 0; i < aLen; i++ {
		if (i > 1 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

/*EncodeTime
 * 字符串转换为time.Time
 */
func EncodeTime(toBeCharge string) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	return theTime
}

//SliToStr 任意类型的数组转换为字符串,插入 "name"，例如：["1", "2", "3", "4"] to "1name2name3name4",
//Param  sl string,params []interface{}
//return str string
func SliToStr(sl string, params ...interface{}) (str string) {
	var paramSlice []string
	for _, param := range params {
		paramSlice = append(paramSlice, param.(string))
	}
	str = strings.Join(paramSlice, sl) // Join 方法第2个参数是 string 而不是 rune
	return
}

//SetTree 递归获取树，根据资源父id，tree指针返回
func SetTree(access []map[string]interface{}, pid string, pTree *map[string]interface{}) {
	i := -1
	var dTree map[string]interface{}
	for k, v := range access {
		dTree = v
		aPid := v["Pid"].(string)
		if aPid == pid && i < k {
			i++
			if (*pTree)["children"] == nil {
				se := make([]map[string]interface{}, 0)
				(*pTree)["children"] = se
			}
			(*pTree)["children"] = append((*pTree)["children"].([]map[string]interface{}), dTree)
			sl := (*pTree)["children"].([]map[string]interface{})
			aId := v["Id"].(string)
			SetTree(access, aId, &sl[i])
		}
	}
}

//GetTree 递归获取树，根据资源父id，tree指针返回
func GetTree(access []map[string]interface{}, pid string) map[string]interface{} {
	dTree := make(map[string]interface{})
	SetTree(access, pid, &dTree)
	return dTree
}

//GetList map 转换为字符串数组
func GetList(pTree map[string]interface{}, pList *[]interface{}) {
	*pList = append(*pList, pTree["Id"])
	if _, ok := pTree["children"]; ok {
		for _, v := range pTree["children"].([]interface{}) {
			GetList(v.(map[string]interface{}), pList)
		}
	}
}

////ShuffleSlice 交换
//func ShuffleSlice(slice []string) {
//	for i := range slice {
//		j := rand.Intn(i + 1)
//		slice[i], slice[j] = slice[j], slice[i]
//	}
//}

func UUID() string {
	var err error
	return uuid.Must(uuid.NewV4(), err).String()
}

// StringInSlice 检查数组slice中是否包含字符串v
func StringInSlice(slice []string, v string) bool {
	for _, item := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// StringsJoin string array join
func StringsJoin(strs ...string) string {
	var str string
	var b bytes.Buffer
	strsLen := len(strs)
	if strsLen == 0 {
		return str
	}
	for i := 0; i < strsLen; i++ {
		b.WriteString(strs[i])
	}
	str = b.String()
	return str

}

func ArrayToStr(arr []interface{}) string {
	result := ""
	for i, v := range arr {
		if v == nil {
			result += "NULL"
		} else {
			result += "'" + v.(string) + "'"
		}
		if i < len(arr)-1 {
			result += ","
		}
	}
	return result
}
func ArrayToString(arr []string) string {
	var result string
	for _, i := range arr { //遍历数组中所有元素追加成string
		result = result + i + "/"
	}
	return result
}
