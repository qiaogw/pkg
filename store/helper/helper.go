package helper

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

func SearchFileList(dirpath, pathStr, searchStr, baseURL, ignore string, filelist []*Entries) []*Entries {
	finfo, _ := ioutil.ReadDir(dirpath)
	var info DirInfo
	info.Path = pathStr
	//var filelist []*seaweed.Entries
	for _, v := range finfo {
		var e Entries
		beego.Debug(v.Name())

		e.Name = v.Name()
		e.FileSize = v.Size()
		e.IsDirectory = v.IsDir()
		e.Mode = v.Mode()
		e.Mtime = v.ModTime()
		e.FullPath = filepath.Join(pathStr, e.Name)
		e.URLDir = baseURL + e.FullPath
		if e.IsDirectory {
			e.Type = "文件夹"
			fullPaht := filepath.Join(dirpath, e.Name)
			beego.Debug(fullPaht)
			filelist = SearchFileList(fullPaht, e.FullPath, searchStr, baseURL, ignore, filelist)
		} else {
			e.Type = path.Ext(e.Name)
			e.Type = e.Type[1:len(e.Type)]
		}
		if path.Ext(e.Name) != ignore && strings.Contains(v.Name(), searchStr) {
			filelist = append(filelist, &e)
		}
	}
	return filelist
}
