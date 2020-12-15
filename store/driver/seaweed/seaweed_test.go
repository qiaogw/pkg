package seaweed

import (

	//"strconv"

	//"path"

	"fmt"

	"github.com/qiaogw/pkg/config"

	//"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeaweedfs(t *testing.T) {
	err := config.LoadConfig()
	r := NewSeaweedfs(`test`)
	fmt.Println(r.baseURL)
	err = r.DeleteDir(`/buckets/buck1/s3data`)
	assert.Nil(t, err)
	//
	//filename := "config.go"
	//fpath := `/js11/test1/1config.go`
	////
	//f, err := os.Open(filename)
	//assert.Nil(t, err)
	//
	//defer f.Close()
	//
	//fi, err := f.Stat()
	//assert.Nil(t, err)
	//
	//_, _, err = r.Put(fpath, f, fi.Size())
	//assert.Nil(t, err)
	//
	//_, err = r.PathInfo(`/js11`)
	//assert.Nil(t, err)
	////fmt.Println(pi)
	//
	//err = r.Rename(fpath, `/github/config.go`)
	//assert.Nil(t, err)
	////
	//fr, _, err := r.Get(fpath)
	//assert.Nil(t, err)
	//assert.NotZero(t, fr)
	//
	//info, err := r.FilerInfo(fpath)
	//assert.Nil(t, err)
	//fmt.Println(info)

	//url := r.FileDir(fpath)
	////fmt.Println(url)
	//
	//url = r.URLDir(fpath)
	//fmt.Println(url)

	//err = r.Delete(fpath)
	//assert.Nil(t, err)
	//
	//err = r.DeleteDir(`/js11/test1`)
	//assert.Nil(t, err)
	//err = r.DeleteDir(`/js11`)
	//assert.Nil(t, err)
}
