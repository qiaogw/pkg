package s3cli

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/qiaogw/log"
	"github.com/qiaogw/pkg/config"
)

var (
	sess *session.Session
	svc  *s3.S3
)

func main() {
	sess = nil
	accessKeyID := config.Config.S3.AccessKeyID
	secretAccessKey := config.Config.S3.SecretAccessKey
	endPoint := config.Config.S3.Endpoint //endpoint设置，不要动
	sess, _ = session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endPoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
	})
	svc = s3.New(sess)
}
func GetAllMyBucket() ([]*s3.Bucket, error) {
	//svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		log.Errorf("Unable to list buckets, %v", err)
		return nil, err
	}
	return result.Buckets, nil
}

func RemoveObject(bucket string, item *s3.Object) (err error) {
	name := aws.StringValue(item.Key)
	dp := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	}
	_, err = svc.DeleteObject(dp)
	if err != nil {
		log.Error("Unable to delete object: ", name, " from bucket :", bucket, err)
	}
	log.Debug("successfully deleted:", aws.StringValue(item.Key))
	return
}

func RestoreObject(bucket string, item *s3.Object, days int64) {
	rparams := &s3.RestoreObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(*item.Key),
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
	_, err := svc.RestoreObject(rparams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeObjectAlreadyInActiveTierError:
				log.Error(s3.ErrCodeObjectAlreadyInActiveTierError, aerr.Error())
			default:
				log.Error(aerr.Error())
			}
		} else {
			log.Error(err)
		}
	}
	return
}

func RemoveObjectsLife() (count int, err error) {
	//svc := s3.New(sess)
	bucket := config.Config.S3.Bucket
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("StoragePath/"),
	}
	//var objkeys []string
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3000)*time.Second)
	defer cancel()
	now := time.Now()
	lifeDay := float64(config.Config.S3.LifeDay - 1)
	err = svc.ListObjectsPagesWithContext(ctx, params, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			t := now.Sub(*content.LastModified).Hours()
			if t > lifeDay*24 {
				if err = RemoveObject(bucket, content); err != nil {
					count++
				}
			}
		}
		return true
	})
	if err != nil {
		log.Error(err)
	}
	return
}

func RestoreObjectsLife(ppath string, days int64) (count int, err error) {
	bucket := config.Config.S3.Bucket
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(ppath),
	}
	//var objkeys []string
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3000)*time.Second)
	defer cancel()
	err = svc.ListObjectsPagesWithContext(ctx, params, func(output *s3.ListObjectsOutput, b bool) bool {
		for _, content := range output.Contents {
			count++
			RestoreObject(bucket, content, days)
		}
		return true
	})
	if err != nil {
		log.Error(err)
	}
	return
}
