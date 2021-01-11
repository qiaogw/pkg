package s3cli

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/astaxie/beego"
	"github.com/qiaogw/pkg/filemanager"
	"github.com/qiaogw/pkg/logs"
	"github.com/qiaogw/pkg/tools"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//Svc s3 管理
type Svc struct {
	svc             *s3.S3
	bucketName      string
	EditableMaxSize int64
	conf            *S3Config
}

// S3Config 配置
type S3Config struct {
	Endpoint          string `label:"地址"`              // 地址
	AccessKeyID       string `label:"AccessKeyID"`     // 地址
	SecretAccessKey   string `label:"SecretAccessKey"` // 地址
	Region            string `label:"对象存储的region"`     // 对象存储的region
	Bucket            string `label:"对象存储的Bucket"`     // 对象存储的Bucket
	Secure            bool   `label:"true代表使用HTTPS"`   // true代表使用HTTPS
	Ignore            string `label:"隐藏文件，S3不支持空目录"`   // 地址
	LifeDay           int64  `label:"存储周期，天"`          // 地址
	DefautRestorePath string `label:"默认恢复文件前缀"`
	TaskTime          string `label:"删除超期文件时间"`
	TempDir           string `label:"临时文件夹"`
	MountDir          string `label:"mount文件夹"`          // 地址
	CacheDir          string `label:"Cache文件夹"`          // 地址
	MountConfigFile   string `label:"mountConfigFile地址"` // 地址
	LogFile           string `label:"LogFile地址"`         // 地址
}

// NewSvc 新的svc
func NewSvc(s3Conf *S3Config) *Svc {
	//sess = nil
	accessKeyID := s3Conf.AccessKeyID
	secretAccessKey := s3Conf.SecretAccessKey
	endPoint := s3Conf.Endpoint //endpoint设置，不要动
	sess, _ := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endPoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
	})
	svc := new(Svc)
	svc.svc = s3.New(sess)
	svc.bucketName = s3Conf.Bucket

	return svc
}

// ListBuckets 查看S3中包含的bucket
func (s *Svc) ListBuckets() ([]*s3.Bucket, error) {
	//svc := s3.New(sess)
	result, err := s.svc.ListBuckets(nil)
	if err != nil {
		logs.Error("Unable to list buckets, %v", err)
		return nil, err
	}
	return result.Buckets, nil
}

// BucketName 查看默认bucket
func (s *Svc) BucketName() string {
	return s.bucketName
}

// CreateBucket 创建bucket
func (s *Svc) CreateBucket(bucketName string) (err error) {
	_, err = s.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		logs.Errorf("Unable to create bucket %q, %v", bucketName, err)
		return err
	}
	err = s.svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	return
	//return s.svc.CreateBucketWithContext(ctx, bucketName, minio.MakeBucketOptions{Region: m.Region})
}

// RemoveBucket Remove bucket
func (s *Svc) RemoveBucket(bucket string) (err error) {
	_, err = s.svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		logs.Errorf("Unable to delete bucket %q, %v", bucket, err)
		return err
	}

	// Wait until bucket is deleted before finishing
	err = s.svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	return
}

//  ListObjects 查看某个bucket中包含的文件/文件夹
func (s *Svc) ListObjects(bucket, prefix string, maxKeys int64, isDelimiter bool) (dirs []os.FileInfo) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
		//MaxKeys:   aws.Int64(maxKeys),
		Marker: aws.String(""),
	}
	if isDelimiter {
		params.Delimiter = aws.String("/")
	}
	var mtime time.Time
	class := ""
	pages := int64(0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3000)*time.Second)
	defer cancel()
	err := s.svc.ListObjectsPagesWithContext(ctx, params, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			if len(prefix) > 0 {
				*content.Key = strings.TrimPrefix(*content.Key, prefix)
			}
			mtime = *content.LastModified
			class = *content.StorageClass
			obj := NewFileInfo(*content)
			dirs = append(dirs, obj)
		}
		pages++
		if maxKeys > 0 {
			return pages <= maxKeys
		}
		return true
	})
	if err != nil {
		beego.Error(err)
		return
	}

	if prefix == "" {
		prefix = "/"
	}
	class = ""
	err = s.GetDir(&dirs, bucket, prefix, class, mtime)
	return
}

// GetDir 获取文件夹树
func (s *Svc) GetDir(dirs *[]os.FileInfo, bucket, prefix, class string, mtime time.Time) (err error) {
	var ld []FilePrefix
	dirfile := filepath.Join(s.conf.TempDir, bucket+".json")
	jdata, _ := ioutil.ReadFile(dirfile)
	err = json.Unmarshal(jdata, &ld)
	if err != nil {
		return
	}
	for _, v := range ld {
		if v.Pid == prefix {
			vs := s3.Object{
				Key:          aws.String(v.Key),
				Size:         aws.Int64(0),
				LastModified: aws.Time(mtime),
				StorageClass: aws.String(class),
			}
			obj := NewFileInfo(vs)
			*dirs = append(*dirs, obj)
		}
	}
	return
}

// GetDirInfo list某个bucket中包含的文件夹
func (s *Svc) GetDirInfo(bucket, pptath string) (dirs []os.FileInfo) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(pptath),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300000)*time.Second)
	defer cancel()
	pageNum := 0
	var ld []FilePrefix
	err := s.svc.ListObjectsPagesWithContext(ctx, params, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			ld = getFilePrefix(path.Dir(*content.Key), pptath, ld)
			pageNum++
			return pageNum <= 50
		}
		return true
	})
	if err != nil {
		return
	}
	nj, _ := json.Marshal(ld)
	dirFile := filepath.Join(s.conf.TempDir, bucket+".json")
	err = tools.WriteFileByte(dirFile, nj, true)
	if err != nil {
		beego.Error(err)
	}
	return
}

//List 列出文件对象包括文件夹
func (s *Svc) List(bucket, ppath string, isDelimiter bool, sortBy ...string) (err error, exit bool, dirs []os.FileInfo) {
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
	dirs = s.ListObjects(bucket, objectPrefix, 0, isDelimiter)
	// dirs = s.GetDirInfo(bucket, objectPrefix)

	if !forceDir && len(dirs) == 0 && err != nil {
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
	return
}

// RemoveObject 删除某个bucket中的对象文件
func (s *Svc) RemoveObject(bucket string, itemName string) (err error) {
	//name := aws.StringValue(item)
	dp := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(itemName),
	}
	_, err = s.svc.DeleteObject(dp)
	if err != nil {
		logs.Error("Unable to delete object ", itemName, "from bucket", bucket, err)
	}
	logs.Info("successfully deleted", itemName)
	return
}

// RemoveObject 恢复某个bucket中的冷存储对象文件
func (s *Svc) RestoreObject(bucket, itemName string, days int64) (err error) {
	rparams := &s3.RestoreObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(itemName),
		RestoreRequest: &s3.RestoreRequest{
			Days: aws.Int64(days),
			GlacierJobParameters: &s3.GlacierJobParameters{
				Tier: aws.String("Expedited"),
				// 取回选项，支持三种取
				//值：[Expedited|Standard|
				//Bulk]。
				//Expedited表示快速取回对
				//象，取回耗时1~5 min，
				//Standard表示标准取回对
				//象，取回耗时3~5 h，
				//Bulk表示批量取回对象，
				//取回耗时5~12 h。
				//默认取值为Standard。
			},
		},
	}
	_, err = s.svc.RestoreObject(rparams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeObjectAlreadyInActiveTierError:
				logs.Error(s3.ErrCodeObjectAlreadyInActiveTierError, aerr.Error())
			default:
				logs.Error(aerr.Error())
			}
		} else {
			logs.Error(err)
		}
	}
	return
}

// RemoveObjectsLife 删除某个Bucket重的超期对象文件
func (s *Svc) RemoveObjectsLife() (count int, err error) {
	logs.Info("开始删除超期文件。。。")
	//svc := s3.New(sess)
	//bucket := s3Conf.Bucket
	params := &s3.ListObjectsInput{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String("StoragePath/"),
	}
	logs.Info(*s.svc.Config.Endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3000)*time.Second)
	defer cancel()
	now := time.Now()
	lifeDay := float64(s.conf.LifeDay - 1)
	err = s.svc.ListObjectsPagesWithContext(ctx, params, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			t := now.Sub(*content.LastModified).Hours()
			if t > lifeDay*24 {
				count++
				s.RemoveObject(s.bucketName, *content.Key)
			}
		}
		return true
	})
	if err != nil {
		logs.Error(err)
	}
	logs.Infof("共删除超期文件 %v 份！！！", count)
	return
}

// RestoreObjectsLife 恢复某个Bucket中的所有冷存储对象文件
func (s *Svc) RestoreObjectsLife(bucketName, ppath string, days int64) (count int, err error) {
	logs.Info("开始恢复文件。。。")
	bucket := bucketName
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(ppath),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3000)*time.Second)
	defer cancel()
	err = s.svc.ListObjectsPagesWithContext(ctx, params, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			if *content.StorageClass == "GLACIER" {
				count++
				s.RestoreObject(bucket, *content.Key, days)
			}
		}
		return true
	})
	if err != nil {
		logs.Error(err)
	}
	logs.Infof("共恢复文件 %v 份！！！", count)
	return
}

// RemoveDir 删除文件夹
func (s *Svc) RemoveDir(prefix string) (err error) {
	objectName := strings.TrimPrefix(prefix, `/`)
	if !strings.HasSuffix(objectName, `/`) {
		objectName += `/`
	}
	if objectName == `/` {
		return s.clear(s.bucketName)
	}
	params := &s3.ListObjectsInput{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(objectName),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3000)*time.Second)
	defer cancel()
	err = s.svc.ListObjectsPagesWithContext(ctx, params, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			if *content.StorageClass == "GLACIER" {
				s.RemoveObject(s.bucketName, *content.Key)
			}
		}
		return true
	})
	return nil
}
func (s *Svc) clear(bucket string) (err error) {
	iter := s3manager.NewDeleteListIterator(s.svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})

	if err := s3manager.NewBatchDeleteWithClient(s.svc).Delete(aws.BackgroundContext(), iter); err != nil {
		logs.Errorf("Unable to delete objects from bucket %q, %v", bucket, err)
		return err
	}
	return
}

// Put 提交数据
func (s *Svc) Put(reader io.Reader, objectName string) (err error) {
	input := &s3.PutObjectInput{
		//ACL:    aws.String("authenticated-read"),
		Body:   aws.ReadSeekCloser(reader),
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	}
	logs.Debug(objectName, s.bucketName)
	_, err = s.svc.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				logs.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logs.Errorf(err.Error())
		}
		return
	}
	return
}

//
// Get 获取数据
func (s *Svc) Get(bucket, ppath string) (io.Reader, error) {
	objectName := strings.TrimPrefix(ppath, `/`)
	beego.Info(objectName)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectName),
	}
	state, err := s.Stat(s.bucketName, ppath)
	if *state.StorageClass == "GLACIER" {
		if state.Restore == nil {
			s.RestoreObject(s.bucketName, ppath, 7)
			return nil, errors.New("对象为冷存储，开始恢复，请在5分钟后重试")
		}
		if strings.Index(*state.Restore, `ongoing-request="false"`) < 0 {
			return nil, errors.New("对象为冷存储，恢复处理中，请在5分钟后重试")
		}
	}
	result, err := s.svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				logs.Error(s3.ErrCodeNoSuchKey, aerr.Error())
			case s3.ErrCodeInvalidObjectState:
				logs.Error(s3.ErrCodeInvalidObjectState, aerr.Error())
			default:
				logs.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logs.Error(err.Error())
		}
		return nil, err
	}

	return result.Body, nil
}

//
// Stat 获取对象信息
func (s *Svc) Stat(bucket, name string) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	}
	result, err := s.svc.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				logs.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logs.Error(err.Error())
		}
		return nil, err
	}
	return result, nil
}

//
// Exists 对象是否存在
func (s *Svc) Exists(ppath string) (bool, error) {
	_, err := s.Stat(s.bucketName, ppath)
	if err != nil {
		return false, err
	}
	return true, err
}

//
