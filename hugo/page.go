package hugo

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/njones/particle"
	"github.com/qiaogw/pkg/tools"
)

func (c *SiteConfig) GetPageList(hugoConfig *Config, menu string) (flist []os.FileInfo, err error) {
	configfile := filepath.Join(hugoConfig.Dir, c.Title, "content", "zh", menu)
	flist, err = ioutil.ReadDir(configfile)
	return
}
func (c *SiteConfig) GetPageInfoList(hugoConfig *Config, menu string) (infoList []FrontMatter, err error) {
	configfile := filepath.Join(hugoConfig.Dir, c.Title, "content", "zh", menu)
	flist, err := ioutil.ReadDir(configfile)
	for _, v := range flist {
		if v.Name() == "_index.md" {
			continue
		}
		fileExt := filepath.Ext(v.Name())

		// 包含.的扩展名
		if fileExt != ".md" {
			continue
		}
		path := filepath.Join(configfile, v.Name())
		info, _, _ := GetFranMatter(path)
		info.Path = filepath.Join("content", "zh", menu, v.Name())
		infoList = append(infoList, info)
	}
	return
}

func (c *SiteConfig) GetPageInfo(hugoConfig *Config, path string) (info FrontMatter, err error) {
	configfile := filepath.Join(hugoConfig.Dir, c.Title, path)
	info, content, _ := GetFranMatter(configfile)
	info.Path = path
	info.Content = string(content)
	return
}

func GetFranMatter(path string) (metadata FrontMatter, content []byte, err error) {
	r, err := os.Open(path)
	defer r.Close()
	//content, err = particle.JSONEncoding.DecodeReader(r, &metadata)
	//beego.Error(err, content, metadata)
	//if err != nil {
	//	content, err = particle.TOMLEncoding.DecodeReader(r, &metadata)
	//}
	//beego.Error(err, content, metadata)
	//if err != nil {
	//	content, err = particle.YAMLEncoding.DecodeReader(r, &metadata)
	//}
	content, err = particle.YAMLEncoding.DecodeReader(r, &metadata)
	return
}
func (c *SiteConfig) SetPageInfo(hugoConfig *Config, page *FrontMatter) (err error) {
	configFile := filepath.Join(hugoConfig.Dir, c.Title, page.Path)
	body := page.Content
	page.Content = ""
	front := particle.YAMLEncoding.EncodeToString([]byte(page.Content), page)
	err = tools.WriteFile(configFile, front, true)
	err = tools.WriteFile(configFile, body, false)
	err = c.BuildSite(hugoConfig)
	return
}
