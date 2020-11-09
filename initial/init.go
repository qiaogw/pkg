// Copyright 2018 cloudy itcloudy@qq.com.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package initial

// var (
// 	// 初始化安装
// 	Installed sql.NullBool
// 	// 安装版本
// 	installedSchemaVer float64
// 	// 安装时间
// 	installedTime time.Time
// 	// 重现加载
// 	// reload bool
// )

// func SetInstalled(lockFile string) error {
// 	now := time.Now()
// 	err := ioutil.WriteFile(lockFile, []byte(now.Format(`2006-01-02 15:04:05`)+"\n"+fmt.Sprint(conf.Config.DBUpdateToVersion)), os.ModePerm)
// 	if err != nil {
// 		return err
// 	}
// 	installedTime = now
// 	installedSchemaVer = conf.Config.DBUpdateToVersion
// 	Installed.Valid = true
// 	Installed.Bool = true
// 	return nil
// }
// func LoadInitData() {
// 	if conf.Config.Init.Enable {
// 		// cfg := conf.Config.DB
// 		conf.GetDBConnection("")
// 		logs.Logger.Info("start load init data")
// 		if conf.Config.Init.API != "" {
// 			logs.Logger.Info("start load init api data")
// 			// initApi()
// 			logs.Logger.Info("load init api data success")

// 		}
// 	} else {
// 		logs.Logger.Error("start load init data failed,config file init enable is false")
// 		os.Exit(-1)
// 	}

// }

// func isInstalled() bool {
// 	dir, _ := os.Getwd()
// 	if !Installed.Valid {
// 		lockFile := filepath.Join(dir, `installed.lock`)
// 		if info, err := os.Stat(lockFile); err == nil && !info.IsDir() {
// 			if b, e := ioutil.ReadFile(lockFile); e == nil {
// 				content := string(b)
// 				content = strings.TrimSpace(content)
// 				lines := strings.Split(content, "\n")
// 				switch len(lines) {
// 				case 2:
// 					installedSchemaVer, _ = strconv.ParseFloat(strings.TrimSpace(lines[1]), 64)
// 					fallthrough
// 				case 1:
// 					installedTime, _ = time.Parse(`2006-01-02 15:04:05`, strings.TrimSpace(lines[0]))
// 				}
// 			}
// 			Installed.Valid = true
// 			Installed.Bool = true
// 		}
// 	}

// 	if Installed.Bool && conf.Config.DBUpdateToVersion > installedSchemaVer {
// 		var upgraded bool
// 		// err := createdb()
// 		// if err == nil {
// 		// 	upgraded = true
// 		// } else {
// 		// 	logs.Panic(`数据库表结构需要升级！`)
// 		// }
// 		if upgraded {
// 			installedSchemaVer = conf.Config.DBUpdateToVersion
// 			ioutil.WriteFile(filepath.Join(dir, `installed.lock`), []byte(installedTime.Format(`2006-01-02 15:04:05`)+"\n"+fmt.Sprint(conf.Config.DBUpdateToVersion)), os.ModePerm)
// 		}
// 	}
// 	return Installed.Bool
// }

// func dbSetup() {
// 	// conf.GetDBConnection("init")
// 	if !isInstalled() {
// 		createdb()
// 	}
// 	conf.SqlxDB.Close()
// 	// conf.SqlxDB.Close()
// }
