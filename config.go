package nuwa

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var configFile = "./config.json"
var privateConfigFile = "./private.json"

type config struct {
	CrtFile       string `json:"crt_file"`
	KeyFile       string `json:"key_file"`
	SSLPort       string `json:"ssl_port"`
	Port          string `json:"port"`
	Host          string `json:"host"`
	CrossDomain   string `json:"cross_domain"`
	PostMaxMemory int64  `json:"post_max_memory"`
	UpdateDir     string `json:"update_dir"`
	UpdatePath    string `json:"update_path"`
}

var _config *config

func SetConfigFile(c string) {
	configFile = c
}

func NewConfig() *config {
	return &config{}
}
func SetConfig(c *config) {
	_config = c
}

func Config() *config {
	if _config == nil {
		_config = &config{
			Host:          "localhost",
			Port:          "8808",
			CrossDomain:   "*",
			PostMaxMemory: 1024 * 1024 * 10,
			UpdateDir:     "",
			UpdatePath:    "",
		}

		_config.ReadConfig(configFile)
		_config.ReadPrivateConfig(privateConfigFile)
	}
	return _config
}

func (c *config) ReadConfig(file string) {
	configText, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("配置文件读取错误,启动默认配置:", err.Error())
		return
	}
	err = json.Unmarshal(configText, c)
	if err != nil {
		log.Println("配置文件错误,启动失败:", err.Error())
		os.Exit(0)
	}
}

func (c *config) ReadPrivateConfig(file string) {
	configText, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("未加载", err.Error())
		return
	}
	err = json.Unmarshal(configText, c)
	if err != nil {
		log.Println("未加载", err.Error())
	}
}
