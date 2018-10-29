package main

import (
	"time"
	log "github.com/sirupsen/logrus"
	"saber/registry"
	"saber/proxy"
	_ "net/http/pprof"
	"net/http"
)

func main() {
	//远程获取pprof数据
	go func() {
		log.Println(http.ListenAndServe("localhost:8081", nil))
	}()

	t := time.Now()
	log.Info("ready to start!")
	//命令行参数

	//runtime.GOMAXPROCS(o.NCPU)
	//读取registry配置
	e := registry.NewRegistry()
	defer e.Close()
	//初始化TODO
	r := &proxy.Router{}

	e.LoadNodes(r)
	e.LoadSlots(r)
	log.Println("slot count :", len(r.Slots))

	//启动服务TODO
	p := proxy.NewProxy(r)
	p.Start(t)
}
