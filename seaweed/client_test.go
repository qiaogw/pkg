package seaweed

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	//"io"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	filerDestination := "http://192.168.0.140:8888/github/"
	//fileOrDirs := []string{"client.go"}
	//fs := make([]io.Reader, 0)
	sw, err := NewSeaweed(filerDestination)
	assert.Nil(t, err)
	//for _, fi := range fileOrDirs {
	//	fmt.Println(fi)
	//	f, err := os.Open(fi)
	//	assert.Nil(t, err)
	//	fs = append(fs, f)
	//}
	//defer f.Close()
	fs, err := os.Open("/Users/qgw/Downloads/ApiPost_3.2.2.dmg")
	assert.Nil(t, err)
	defer fs.Close()
	dispath := `/github/`
	err = sw.Put1(fs, dispath)
	assert.Nil(t, err)
}
