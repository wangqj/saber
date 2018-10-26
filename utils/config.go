package utils

import (
	"github.com/koding/multiconfig"
	"os"
	log "github.com/sirupsen/logrus"
)

type Option struct {
	Port         int `default:"16379"`
	NCPU         int
	RegistryAdrr string
}

//读取本地配置文件TODO

func LoadConf() (*Option) {
	m := multiconfig.NewWithPath("config.toml") // supports TOML, JSON and YAML

	serverConf := new(Option)
	err := m.Load(serverConf)
	if err != nil {
		log.Errorln(err)
		os.Exit(0)
	}
	m.MustLoad(serverConf) // Panic's if there is any error
	log.Println("load config =", serverConf)
	return serverConf
}
