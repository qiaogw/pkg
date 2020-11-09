package caddy

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//
//func (c *SiteConfig) NewSite(hugoConfig *Config) (err error) {
func GetCaddy() *Config {
	return DefaultConfig
}
func (c *Config) NewSite(conf, title string) (err error) {
	confPath := filepath.Dir(c.Caddyfile)
	siteConfPath := filepath.Join(confPath, "caddy", title+".conf")
	err = ioutil.WriteFile(siteConfPath, []byte(conf), os.ModePerm)
	if err == nil {
		err = c.Restart()
	}
	return
}

func (c *Config) RemoveSite(title string) error {
	confPath := filepath.Dir(c.Caddyfile)
	siteConfPath := filepath.Join(confPath, "caddy", title+".conf")
	err := os.Remove(siteConfPath)
	if os.IsNotExist(err) {
		err = nil
	}
	err = c.Restart()
	return err
}
