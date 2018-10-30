package utils

import (
	"github.com/koding/multiconfig"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Option struct {
	Port          int `default:"16379"`
	NCPU          int
	RegistryAdrr  string
	SessionBuffer int
	DataBuffer    int
}

var option *Option
var once sync.Once
//读取本地配置文件TODO

func GetConf() (*Option) {
	GetConfByPath("config.toml")
	return option
}

func GetConfByPath(path string) (*Option) {
	once.Do(func() {
		m := multiconfig.NewWithPath(path) // supports TOML, JSON and YAML

		option = new(Option)
		err := m.Load(option)
		if err != nil {
			log.Errorln(err)
			//os.Exit(0)
		}
		m.MustLoad(option) // Panic's if there is any error
		log.Println("load config =", option)
	})
	return option
}
