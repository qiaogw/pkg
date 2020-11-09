package helper

import (
	//"github.com/qiaogw/pkg/store"
	"regexp"
	"strings"
)

var (
	temporaryFileRegexp  *regexp.Regexp
	persistentFileRegexp *regexp.Regexp
	anyFileRegexp        *regexp.Regexp
	placeholderRegexp    = regexp.MustCompile(`\[storage:[\d]+\]`)
)

const (
	DefaultUploadURLPath = `/public/upload/`
	DefaultUploadDir     = `./public/upload`
)

var (
	// UploadURLPath 上传文件网址访问路径
	UploadURLPath = DefaultUploadURLPath

	// UploadDir 定义上传目录（首尾必须带“/”）
	UploadDir = DefaultUploadDir

	// AllowedUploadFileExtensions 被允许上传的文件的扩展名
	AllowedUploadFileExtensions = []string{
		`.jpeg`, `.jpg`, `.gif`, `.png`,
	}

	// FileTypeIcon 文件类型icon
	//FileTypeIcon = store.FileTypeIcon

	// DetectFileType 根据文件扩展名判断文件类型
	//DetectFileType = store.DetectType

	// TypeRegister 注册文件扩展名
	//TypeRegister = store.TypeRegister
)

// URLToFile 文件网址转为存储路径
func URLToFile(fileURL string) string {
	filePath := strings.TrimPrefix(fileURL, UploadURLPath)
	filePath = strings.TrimSuffix(UploadDir, `/`) + `/` + strings.TrimPrefix(filePath, `/`)
	return filePath
}

func ExtensionRegister(extensions ...string) {
	AllowedUploadFileExtensions = append(AllowedUploadFileExtensions, extensions...)
}

func ExtensionRegexpEnd() string {
	extensions := make([]string, len(AllowedUploadFileExtensions))
	for index, extension := range AllowedUploadFileExtensions {
		extensions[index] = regexp.QuoteMeta(extension)
	}
	return `(` + strings.Join(extensions, `|`) + `)`
}

func init() {
	Init()
}

func Init() {
	ruleEnd := ExtensionRegexpEnd()
	temporaryFileRegexp = regexp.MustCompile(UploadURLPath + `[\w-]+/0/[\w]+` + ruleEnd)
	persistentFileRegexp = regexp.MustCompile(UploadURLPath + `[\w-]+/([^0]|[0-9]{2,})/[\w]+` + ruleEnd)
	anyFileRegexp = regexp.MustCompile(UploadURLPath + `[\w-]+/([\w-]+/)+[\w-]+` + ruleEnd)
}

// ParseTemporaryFileName 从文本中解析出临时文件名称
var ParseTemporaryFileName = func(s string) []string {
	files := temporaryFileRegexp.FindAllString(s, -1)
	return files
}

// ParsePersistentFileName 从文本中解析出正式文件名称
var ParsePersistentFileName = func(s string) []string {
	files := persistentFileRegexp.FindAllString(s, -1)
	return files
}

// ParseAnyFileName 从文本中解析出任意上传文件名称
var ParseAnyFileName = func(s string) []string {
	files := anyFileRegexp.FindAllString(s, -1)
	return files
}

// ReplaceAnyFileName 从文本中替换任意上传文件名称
var ReplaceAnyFileName = func(s string, repl func(string) string) string {
	return anyFileRegexp.ReplaceAllStringFunc(s, repl)
}

// ReplacePlaceholder 从文本中替换占位符
var ReplacePlaceholder = func(s string, repl func(string) string) string {
	return placeholderRegexp.ReplaceAllStringFunc(s, func(find string) string {
		id := find[9 : len(find)-1]
		return repl(id)
	})
}
