package s3

import (
	"fmt"
	"log"
	"testing"

	"github.com/qiaogw/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestS3(t *testing.T) {
	err := config.LoadConfig("/Users/qgw/proj/github.com/qiaogw//conf/app.toml")
	assert.Nil(t, err)
	fmt.Println("config.Config.S3.Endpoint:", config.Config.S3.Endpoint)
	r, err := New(1024)
	assert.Nil(t, err)
	log.Printf("%#v\n", r.client) // minioClient is now setup
	fmt.Println("config.Config.S3.bucketName:", r.bucketName)
	err = r.Mkbucket("444444")
	assert.Nil(t, err)
	err = r.RemoveBucket("444444")
	assert.Nil(t, err)
	b, err := r.ListBuckets()
	fmt.Println("config.Config.S3.ListBuckets:", b)
	assert.Nil(t, err)
	err = r.RemoveBucket("3333333")
	assert.Nil(t, err)
	b, err = r.ListBuckets()
	fmt.Println("config.Config.S3.ListBuckets:", b)
	assert.Nil(t, err)
}
