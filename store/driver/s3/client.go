package s3

import (
	"context"
	"path/filepath"

	"github.com/qiaogw/pkg/tools"

	//"errors"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"github.com/qiaogw/pkg/filemanager"

	//"crypto/tls"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	//"net/http"
	"github.com/qiaogw/pkg/config"
)

var (
//m = config.Config.S3
)

type S3Manager struct {
	client          *minio.Client
	bucketName      string
	EditableMaxSize int64
}

func New(editableMaxSize int64) (*S3Manager, error) {
	m := config.Config.S3
	client, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKeyID, m.SecretAccessKey, ""),
		Secure: m.Secure,
		Region: m.Region,
	})
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	mgr := &S3Manager{
		client:          client,
		bucketName:      m.Bucket,
		EditableMaxSize: editableMaxSize,
	}
	return mgr, err
}
func (s *S3Manager) BucketName() string {
	return s.bucketName
}
func (s *S3Manager) Mkbucket(bucketName string) error {
	m := config.Config.S3
	ctx := context.Background()
	return s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: m.Region})
}
func (s *S3Manager) ListBuckets() ([]minio.BucketInfo, error) {
	ctx := context.Background()
	buckets, err := s.client.ListBuckets(ctx)
	return buckets, err
}
func (s *S3Manager) RemoveBucket(bucketName string) error {
	ctx := context.Background()
	err := s.client.RemoveBucket(ctx, bucketName)
	return err
}

func (s *S3Manager) Mkdir(ppath, newName string) error {
	ctx := context.Background()
	objectName := strings.TrimPrefix(ppath, `/`)
	objectName = path.Join(objectName, newName)
	if !strings.HasSuffix(objectName, `/`) {
		objectName += `/`
	}
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, nil, 0, minio.PutObjectOptions{})
	return err
}

func (s *S3Manager) Rename(ppath, newName string) error {
	ctx := context.Background()
	objectName := strings.TrimPrefix(ppath, `/`)
	// Source object
	src := minio.CopySrcOptions{
		Bucket: s.bucketName,
		Object: objectName,
	}
	//(s.bucketName, objectName, nil)
	newName = strings.TrimPrefix(newName, `/`)
	dst := minio.CopyDestOptions{
		Bucket: s.bucketName,
		Object: newName,
	}
	//(s.bucketName, newName, nil, nil)
	//if err != nil {
	//	return err
	//}

	// Initiate copy object.
	_, err := s.client.CopyObject(ctx, dst, src)
	if err != nil {
		return err
	}
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
		VersionID:        "myversionid",
	}
	err = s.client.RemoveObject(ctx, s.bucketName, objectName, opts)
	return err
}

func (s *S3Manager) Chown(ppath string, uid, gid int) error {
	return nil
}

func (s *S3Manager) Chmod(ppath string, mode os.FileMode) error {
	return nil
}

func (s *S3Manager) Search(ppath string, prefix string, num int) []string {
	//ctx := context.Background()
	var paths []string
	//doneCh := make(chan struct{})
	//defer close(doneCh)
	objectPrefix := path.Join(ppath, prefix)
	objectPrefix = strings.TrimPrefix(objectPrefix, `/`)
	//objectCh := s.client.ListObjects(ctx, s.bucketName, objectPrefix, false, doneCh)
	//for object := range objectCh {
	//	if object.Err != nil {
	//		continue
	//	}
	//	paths = append(paths, object.Key)
	//}
	//return paths
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix:    objectPrefix,
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			beego.Error(object.Err)
			return nil
		}
		paths = append(paths, object.Key)
	}
	return paths
}

func (s *S3Manager) Remove(ppath string) error {
	if len(ppath) == 0 {
		return errors.New("path invalid")
	}
	if strings.HasSuffix(ppath, `/`) {
		return s.RemoveDir(ppath)
	}
	objectName := strings.TrimPrefix(ppath, `/`)
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
		VersionID:        "myversionid",
	}
	return s.client.RemoveObject(context.Background(), s.bucketName, objectName, opts)
}

func (s *S3Manager) RemoveDir(ppath string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	if !strings.HasSuffix(objectName, `/`) {
		objectName += `/`
	}
	if objectName == `/` {
		return s.Clear()
	}
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
		VersionID:        "myversionid",
	}
	s.client.RemoveObject(context.Background(), s.bucketName, objectName, opts)
	doneCh := make(chan struct{})
	defer close(doneCh)
	listOpts := minio.ListObjectsOptions{
		Prefix:    objectName,
		Recursive: true,
	}
	objectCh := s.client.ListObjects(context.Background(), s.bucketName, listOpts)
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		if len(object.Key) == 0 {
			continue
		}
		err := s.client.RemoveObject(context.Background(), s.bucketName, object.Key, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

// Clear 清空所有数据【慎用】
func (s *S3Manager) Clear() error {
	deleted := make(chan minio.ObjectInfo)
	defer close(deleted)
	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	removeObjects := s.client.RemoveObjects(context.Background(), s.bucketName, deleted, opts)
	for removeObject := range removeObjects {
		if removeObject.Err != nil {
			return removeObject.Err
		}
	}
	return nil
}

//func (s *S3Manager) Upload(ctx echo.Context, ppath string) error {
//	fileSrc, fileHdr, err := ctx.Request().FormFile(`file`)
//	if err != nil {
//		return err
//	}
//	defer fileSrc.Close()
//	objectName := path.Join(ppath, fileHdr.Filename)
//	return s.Put(fileSrc, objectName, fileHdr.Size)
//}

// Put 提交数据
func (s *S3Manager) Put(reader io.Reader, objectName string, size int64) (err error) {
	opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}
	objectName = strings.TrimPrefix(objectName, `/`)
	_, err = s.client.PutObject(context.Background(), s.bucketName, objectName, reader, size, opts)
	return
}

// Get 获取数据
func (s *S3Manager) Get(ppath string) (*minio.Object, error) {
	objectName := strings.TrimPrefix(ppath, `/`)
	f, err := s.client.GetObject(context.Background(), s.bucketName, objectName, minio.GetObjectOptions{})
	//if err != nil {
	//	return f, errors.WithMessage(err, objectName)
	//}
	return f, err
}

// Stat 获取对象信息
func (s *S3Manager) Stat(ppath string) (minio.ObjectInfo, error) {
	objectName := strings.TrimPrefix(ppath, `/`)
	f, err := s.client.StatObject(context.Background(), s.bucketName, objectName, minio.StatObjectOptions{})
	//if err != nil {
	//	return f, errors.WithMessage(err, objectName)
	//}
	return f, err
}

// Exists 对象是否存在
func (s *S3Manager) Exists(ppath string) (bool, error) {
	_, err := s.Stat(ppath)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case minio.ErrorResponse:
			return v.StatusCode != http.StatusNotFound, nil
		}
		return false, err
	}
	return true, err
}

//
//func (s *S3Manager) Download(ctx echo.Context, ppath string) error {
//	f, err := s.Get(ppath)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	fileName := path.Base(ppath)
//	inline := ctx.Formx(`inline`).Bool()
//	return ctx.Attachment(f, fileName, inline)
//}

func (s *S3Manager) List(ppath string, sortBy ...string) (err error, exit bool, dirs []os.FileInfo) {
	ctx := context.Background()
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectPrefix := strings.TrimPrefix(ppath, `/`)
	words := len(objectPrefix)
	var forceDir bool
	if words == 0 {
		forceDir = true
	} else {
		if strings.HasSuffix(objectPrefix, `/`) {
			forceDir = true
		} else {
			objectPrefix += `/`
		}
	}
	listOpts := minio.ListObjectsOptions{
		Prefix:    objectPrefix,
		Recursive: true,
	}
	objectCh := s.client.ListObjects(ctx, s.bucketName, listOpts)
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		if len(objectPrefix) > 0 {
			object.Key = strings.TrimPrefix(object.Key, objectPrefix)
		}
		if len(object.Key) == 0 {
			continue
		}
		obj := NewFileInfo(object)
		dirs = append(dirs, obj)
	}
	if !forceDir && len(dirs) == 0 {
		return
	}
	if len(sortBy) > 0 {
		switch sortBy[0] {
		case `time`:
			sort.Sort(filemanager.SortByModTime(dirs))
		case `-time`:
			sort.Sort(filemanager.SortByModTimeDesc(dirs))
		case `name`:
		case `-name`:
			sort.Sort(filemanager.SortByNameDesc(dirs))
		case `type`:
			fallthrough
		default:
			sort.Sort(filemanager.SortByFileType(dirs))
		}
	} else {
		sort.Sort(filemanager.SortByFileType(dirs))
	}
	//if ctx.Format() == "json" {
	//	dirList, fileList := s.ListTransfer(dirs)
	//	data := ctx.Data()
	//	data.SetData(echo.H{
	//		`dirList`:  dirList,
	//		`fileList`: fileList,
	//	})
	//	return ctx.JSON(data), true, nil
	//}
	return
}

func (s *S3Manager) ListTree(dirpath string) []tools.DirBody {
	var allFile []tools.DirBody
	_, _, dirs := s.List(dirpath)
	for _, x := range dirs {
		var chiledren tools.DirBody
		if x.IsDir() {
			realPath := filepath.Join(dirpath, x.Name())
			chiledren.Label = x.Name()
			chiledren.Dir = realPath
			chiledren.Icon = "el-icon-folder"
			chiledren.Children = append(chiledren.Children, s.ListTree(realPath)...)
			//beego.Debug(chiledrens)
			allFile = append(allFile, chiledren)
		}
	}
	return allFile
}

func (s *S3Manager) ListObj(ppath string) (err error, exit bool, dirs []os.FileInfo) {
	ctx := context.Background()
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectPrefix := strings.TrimPrefix(ppath, `StoragePath/`)
	words := len(objectPrefix)
	var forceDir bool
	if words == 0 {
		forceDir = true
	} else {
		if strings.HasSuffix(objectPrefix, `/`) {
			forceDir = true
		} else {
			objectPrefix += `/`
		}
	}
	listOpts := minio.ListObjectsOptions{
		Prefix:    objectPrefix,
		Recursive: true,
	}
	beego.Debug(s.bucketName)
	lists, err := s.client.ListBuckets(ctx)

	if err != nil {
		beego.Debug(err)
	}

	for _, list := range lists {
		beego.Debug(list.Name)
	}
	objectCh := s.client.ListObjects(ctx, s.bucketName, listOpts)
	for object := range objectCh {
		if object.Err != nil {
			beego.Error(object.Err, object)
			continue
		}
		//if len(objectPrefix) > 0 {
		//	object.Key = strings.TrimPrefix(object.Key, objectPrefix)
		//}
		if len(object.Key) == 0 {
			continue
		}
		beego.Debug(object.Key)
		obj := NewFileInfo(object)
		dirs = append(dirs, obj)
	}
	if !forceDir && len(dirs) == 0 {
		return
	}

	//if ctx.Format() == "json" {
	//	dirList, fileList := s.ListTransfer(dirs)
	//	data := ctx.Data()
	//	data.SetData(echo.H{
	//		`dirList`:  dirList,
	//		`fileList`: fileList,
	//	})
	//	return ctx.JSON(data), true, nil
	//}
	return
}
