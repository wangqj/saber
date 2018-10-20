package main

import (
	"time"
	log "github.com/sirupsen/logrus"
	"saber/utils"
	"saber/registry"
	"saber/proxy"
	"runtime"
	_ "net/http/pprof"
	"net/http"
)

func main() {
	//远程获取pprof数据
	go func() {
		log.Println(http.ListenAndServe("localhost:8081", nil))
	}()

	runtime.GOMAXPROCS(8)
	t := time.Now()
	log.Info("ready to start!")
	//命令行参数

	//读取配置文件，校验参数TODO
	//ip port
	o := utils.LoadConf()
	//读取registry配置
	e := registry.NewRegistry(o)
	defer e.Close()
	//初始化TODO
	r := proxy.Router{}

	e.LoadNodes(&r)
	e.LoadSlots(&r)
	log.Println("slot count :", len(r.Slots))

	//启动服务TODO
	p := proxy.NewProxy(o, &r)
	p.Start(t)

}
