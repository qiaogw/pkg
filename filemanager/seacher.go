package filemanager

import (
	//"github.com/webx-top/echo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Search 自动完成查询文件
func Search(prefix string, nums ...int) []string {
	var paths []string
	root := filepath.Dir(prefix)
	if len(root) == 0 {
		root, _ = os.Getwd()
		root = root + "/"
	}
	num := 10
	if len(nums) > 0 {
		num = nums[0]
	}
	dir, _ := ioutil.ReadDir(root)
	for _, d := range dir {
		if len(paths) >= num {
			break
		}
		path := filepath.Join(root, d.Name())
		if d.IsDir() {
			path += "/"
		}
		if len(prefix) == 0 {
			paths = append(paths, path)
			continue
		}
		if strings.HasPrefix(path, prefix) {
			paths = append(paths, path)
			continue
		}
	}
	return paths
}

// SearchDir 自动完成查询文件
func SearchDir(prefix string, nums ...int) []string {
	var paths []string
	root := filepath.Dir(prefix)
	if len(root) == 0 {
		root, _ = os.Getwd()
		root = root + "/"
	}
	num := 10
	if len(nums) > 0 {
		num = nums[0]
	}
	dir, _ := ioutil.ReadDir(root)
	for _, d := range dir {
		if len(paths) >= num {
			break
		}
		if !d.IsDir() {
			continue
		}
		path := filepath.Join(root, d.Name())
		path += "/"
		if len(prefix) == 0 {
			paths = append(paths, path)
			continue
		}
		if strings.HasPrefix(path, prefix) {
			paths = append(paths, path)
			continue
		}
	}
	return paths
}

// SearchFile 自动完成查询文件
func SearchFile(prefix string, nums ...int) []string {
	var paths []string
	root := filepath.Dir(prefix)
	if len(root) == 0 {
		root, _ = os.Getwd()
		root = root + "/"
	}
	num := 10
	if len(nums) > 0 {
		num = nums[0]
	}
	dir, _ := ioutil.ReadDir(root)
	for _, d := range dir {
		if len(paths) >= num {
			break
		}
		if d.IsDir() {
			continue
		}
		path := filepath.Join(root, d.Name())
		if len(prefix) == 0 {
			paths = append(paths, path)
			continue
		}
		if strings.HasPrefix(path, prefix) {
			paths = append(paths, path)
			continue
		}
	}
	return paths
}
