package s3cli

import (
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qiaogw/pkg/tools"
	"github.com/wxnacy/wgo/arrays"
)

func NewFileInfo(objectInfo s3.Object) os.FileInfo {
	return &fileInfo{objectInfo: objectInfo}
}

type fileInfo struct {
	objectInfo s3.Object
}

func (f *fileInfo) Name() string {
	return *f.objectInfo.Key
}

func (f *fileInfo) Size() int64 {
	return *f.objectInfo.Size
}

func (f *fileInfo) Mode() os.FileMode {
	return 0
}

func (f *fileInfo) ModTime() time.Time {
	return *f.objectInfo.LastModified
}

func (f *fileInfo) IsDir() bool {
	return strings.HasSuffix(f.Name(), "/")
}

func (f *fileInfo) Sys() interface{} {
	return f.objectInfo
}

type FilePrefix struct {
	Key  string
	Pid  string
	Time time.Time
}
//
func getFilePrefix(dirctory, prefix string, ld []FilePrefix) []FilePrefix {
	s := strings.Split(dirctory, "/")
	for i, v := range s {
		var td FilePrefix
		if i == 0 {
			td.Pid = "/"
		} else {
			td.Pid = tools.ArrayToString(s[0:i])
		}
		td.Key = v + "/"

		if arrays.Contains(ld, td) < 0 {
			ld = append(ld, td)
		}
	}
	return ld
}
