package main

import (
	_ "net/http/pprof"
	"saber/cmd/proxy/app"
	"net/http"
	log "github.com/sirupsen/logrus"
)



func main() {
	//远程获取pprof数据
	go func() {
		log.Println(http.ListenAndServe("localhost:8081", nil))
	}()
	//runtime.GOMAXPROCS(o.NCPU)
	app.Execute()
}
