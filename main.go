package main

import (
	"time"
	log "github.com/sirupsen/logrus"
	"saber/utils"
	"saber/registry"
	"saber/proxy"
)

func main() {
	t1 := time.Now()
	log.Info("ready to start!")
	//命令行参数

	//读取配置文件，校验参数TODO
	//ip port
	o := utils.LoadConf()
	//读取registry配置
	registry.LoadNodes()
	registry.LoadSlots()
	//初始化TODO

	//启动服务TODO
	proxy.Proxy{}.Start(o)

	t2 := time.Since(t1)
	log.WithFields(log.Fields{"Spend time": t2,}).Info("Successful startup!")
}
