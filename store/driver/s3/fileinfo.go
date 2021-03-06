package s3

import (
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/minio/minio-go/v7"
	//minio "github.com/minio/minio-go"
)

func NewFileInfo(objectInfo minio.ObjectInfo) os.FileInfo {
	return &fileInfo{objectInfo: objectInfo}
}

type fileInfo struct {
	objectInfo minio.ObjectInfo
}

func (f *fileInfo) Name() string {
	pop := strings.LastIndex(f.objectInfo.Key, "/")
	beego.Debug(f.objectInfo.Key, pop)
	name := f.objectInfo.Key
	if pop > -1 {
		name = f.objectInfo.Key[pop:]
	}

	return name
}

func (f *fileInfo) Size() int64 {
	return f.objectInfo.Size
}

func (f *fileInfo) Mode() os.FileMode {
	return 0
}

func (f *fileInfo) ModTime() time.Time {
	return f.objectInfo.LastModified
}

func (f *fileInfo) IsDir() bool {
	return strings.HasSuffix(f.Name(), "/")
}

func (f *fileInfo) Sys() interface{} {
	return f.objectInfo
}
