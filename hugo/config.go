package hugo

import "time"

type Config struct {
	Dir      string `json:"dir"`      //hugo 站点文件目录
	Title    string `json:"title"`    //站点名称
	ThemeDir string `json:"themedir"` //站点名称
}

type SiteConfig struct {
	BaseURL                string `json:"baseURL"`                //站点地址
	Title                  string `json:"title"`                  //站点名称
	Theme                  string `json:"theme"`                  //站点模版
	Paginate               string `json:"paginate"`               //分页开始数量
	SummaryLength          string `json:"summaryLength"`          //摘要字数
	DisqusShortname        string `json:"disqusShortname"`        //评论服务
	DisableLanguages       string `json:"disableLanguages"`       //禁用语言
	DefaultContentLanguage string `json:"defaultContentLanguage"` //默认语言
	Menu                   Menu
}

type Menu struct {
	Main   []MenuConfig
	Footer []MenuConfig
}

type MenuConfig struct {
	Name   string `json:"name"`   //站点地址
	URL    string `json:"URL"`    //站点名称
	Weight int    `json:"weight"` //站点模版
}

type FrontMatter struct {
	Title       string    `json:"title" label:"标题"`       //标题
	Description string    `json:"description" label:"简介"` //页面简介
	Type        string    `json:"type"`                   //页面简介
	Date        time.Time `json:"date" label:"发布日期"`      //发布日期
	Draft       bool      `json:"draft" label:"是否草稿"`     //是否草稿
	Author      string    `json:"author" label:"作者"`      //作者
	Categories  []string  `json:"categories" label:"分类"`  //分类
	Tags        []string  `json:"tags" label:"标签"`        //标签
	Path        string    `json:"path" label:"文件名"`       //标签
	Content     string    `json:"content"`                //文件内容
	SiteId      int64     `json:"siteId"`                 //站点id
	Image       string    `json:"image"`                  //站点id
	BgImage     string    `json:"bg_image"`               //站点id
}
