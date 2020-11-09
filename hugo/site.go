package hugo

import (
	"os"
	"path/filepath"

	"github.com/go-cmd/cmd"
	"github.com/naoina/toml"
	"github.com/pkg/errors"
	"github.com/qiaogw/pkg/logs"
	"github.com/qiaogw/pkg/tools"
)

//var (
//	hugoConfig = config.Config.Hugo
//)

//func GetHugoConfig() *Config {
//	return &hugoConfig
//}

func NewHugo() *SiteConfig {
	return &SiteConfig{
		Paginate:               "4",
		SummaryLength:          "30",
		DefaultContentLanguage: "zh",
	}
}
func (c *SiteConfig) NewSite(hugoConfig *Config) (err error) {

	path := filepath.Join(hugoConfig.Dir, c.Title)
	envCmd := cmd.NewCmd("hugo", "new", "site", path)
	status := <-envCmd.Start()
	// Print each line of STDOUT from Cmd
	for _, line := range status.Stdout {
		logs.Debug(line)
	}
	themeDest := filepath.Join(path, "themes", c.Theme)
	themeSrc := filepath.Join(hugoConfig.ThemeDir, c.Theme)
	siteSrc := filepath.Join(themeSrc, "exampleSite")
	err = tools.CopyDir(siteSrc, path)
	if err != nil {
		logs.Error(err)
		return
	}
	err = tools.CopyDir(themeSrc, themeDest)
	if err != nil {
		logs.Error(err)
		return
	}
	conf := make(map[string]interface{})
	configFile := filepath.Join(hugoConfig.Dir, c.Title, "config.toml")
	err = tools.GetConfigFromPath(configFile, &conf)

	//err = json.Unmarshal(confValue, conf)
	conf["baseurl"] = c.BaseURL
	conf["title"] = c.Title
	conf["theme"] = c.Theme
	err = SaveConfig(configFile, conf)
	if err != nil {
		logs.Error(err)
		return
	}
	err = c.BuildSite(hugoConfig)
	if err != nil {
		logs.Error(err)
	}
	return
}

func (c *SiteConfig) GetConfigFile(hugoConfig *Config) (str string, err error) {
	configfile := filepath.Join(hugoConfig.Dir, c.Title, "config.toml")
	str, err = tools.ReadFile(configfile)
	return
}

//func (c *SiteConfig) GetConfigFileParams(hugoConfig *Config) (fileInfo interface{}, err error) {
//	configfile := filepath.Join(hugoConfig.Dir, c.Title, "config.toml")
//	//fileInfo, err = config.GetConfigFromPath(configfile)
//	return
//}
func (c *SiteConfig) SetConfigFile(hugoConfig *Config, str string) (err error) {
	configFile := filepath.Join(hugoConfig.Dir, c.Title, "config.toml")
	err = tools.WriteFile(configFile, str, true)
	if err == nil {
		conf := make(map[string]interface{})
		err = tools.GetConfigFromPath(configFile, &conf)
		conf["baseurl"] = c.BaseURL
		conf["title"] = c.Title
		conf["theme"] = c.Theme
		err = SaveConfig(configFile, conf)
		if err == nil {
			err = c.BuildSite(hugoConfig)
		}

	}
	return
}
func (c *SiteConfig) BuildSite(hugoConfig *Config) (err error) {
	logs.Debug("BuildSite")
	path := filepath.Join(hugoConfig.Dir, c.Title)
	envCmd := cmd.NewCmd("hugo")
	envCmd.Dir = path
	logs.Debug("BuildSite", path)
	status := <-envCmd.Start()
	// Print each line of STDOUT from Cmd
	for _, line := range status.Stdout {
		logs.Debug(line)
	}
	return status.Error
}

// SaveConfig save global parameters to configFile
func SaveConfig(path string, config interface{}) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "creating dir %s", dir)
		}
	}

	cf, err := os.Create(path)
	if err != nil {
		//log.WithFields(log.Fields{"type": consts.IOError, "error": err}).Error("Create config file failed")
		return err
	}
	defer cf.Close()

	err = toml.NewEncoder(cf).Encode(config)
	if err != nil {
		return err
	}
	return nil
}
