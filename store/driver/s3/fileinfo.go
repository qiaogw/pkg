package s3

import (
	"github.com/minio/minio-go/v7"
	"os"
	"strings"
	"time"
	//minio "github.com/minio/minio-go"
)

func NewFileInfo(objectInfo minio.ObjectInfo) os.FileInfo {
	return &fileInfo{objectInfo: objectInfo}
}

type fileInfo struct {
	objectInfo minio.ObjectInfo
}

func (f *fileInfo) Name() string {
	return f.objectInfo.Key
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
