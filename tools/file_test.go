package tools

import (
	"fmt"
	"os"
	"path/filepath"

	//"regexp"
	"testing"
)

//TestCheckPath 检查文件夹是否存在，不存在则创建
func TestCheckPath(t *testing.T) {
	dir := filepath.Dir("logPath")
	err := CheckPath(dir)
	if err != nil {
		t.Fatal(err)
	}
}

// GetFileSuffix 获取文件扩展名
func TestGetFileSuffix(t *testing.T) {
	dir := filepath.Dir("file.go")
	suf := GetFileSuffix(dir)
	fmt.Println(dir, suf)
	if suf != "." {
		t.Fatal(suf)
	}
}

// LockOrDie 加锁文件夹
func TestLockOrDie(t *testing.T) {
	dir := "./test"
	file := dir + "/" + "test.txt"

	err := CheckPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(file)
	if err != nil {
		t.Fatal(err)
	}
	lf := LockOrDie(dir)
	err = os.Remove(file)
	if err != nil {
		t.Errorf("lock succes ,info  is %v", err)
	}

	defer f.Close()
	defer lf.Unlock()
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
}

// MakeDirectory 创建 directory if is not exists
func TestMakeDirectory(t *testing.T) {

}
