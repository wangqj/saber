package utils

import (
	"github.com/koding/multiconfig"
	"os"
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
	once.Do(func() {
	m := multiconfig.NewWithPath("config.toml") // supports TOML, JSON and YAML

		option = new(Option)
		err := m.Load(option)
	if err != nil {
		log.Errorln(err)
		os.Exit(0)
	}
		m.MustLoad(option) // Panic's if there is any error
		log.Println("load config =", option)
	})
	return option
}
