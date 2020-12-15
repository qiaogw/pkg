// Copyright 2018 cloudy itcloudy@qq.com.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package config

import (
	"github.com/jmoiron/sqlx"
)

type Pagination struct {
	Total   int `json:"total"`
	Size    int `json:"pageSize"`
	Current int `json:"current"`
}
type ResponseJson struct {
	Code    uint        `yaml:"code" json:"code"`       // response code
	Data    interface{} `yaml:"data" json:"data"`       // response data
	Message string      `yaml:"message" json:"message"` // response message
}

var (
	//DBConn *gorm.DB
	//Gdb    *gorm.DB
	//Xdb    *xorm.Engine
	// engine, err = xorm.NewEngine("mysql", "root:123@/test?charset=utf8")

	SqlxDB *sqlx.DB

	// Config global parameters
	Config GlobalConfig
	// Elasticsearch client
	//ElasticClient *elastic.Client
	// casbin
	//Enforcer *casbin.Enforcer
	// 全局语言对象
	//I18nBundles = make(map[string]*i18n.Bundle)
)
