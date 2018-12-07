package main

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"net/http"
	"time"
	"flag"
	"saber/utils"
	"saber/registry"
	_ "net/http/pprof"
	log "github.com/sirupsen/logrus"
	"saber/proxy"
)

var rootCmd = &cobra.Command{
	Use:   "saber",
	Short: "saber is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		//远程获取pprof数据
		go func() {
			log.Println(http.ListenAndServe("localhost:8081", nil))
		}()

		t := time.Now()
		log.Info("ready to start!")
		//命令行参数

		//runtime.GOMAXPROCS(o.NCPU)
		var operation string
		flag.StringVar(&operation, "c", "config.toml", "config file path")
		log.Info("path=", operation)

		utils.GetConfByPath("config.toml")

		//读取registry配置
		e := registry.NewRegistry()
		defer e.Close()
		//初始化
		r := &proxy.Router{}
		e.LoadAll(r)
		log.Println("node count : ", len(r.Nodes), "; slot count :", len(r.Slots))

		//启动服务
		p := proxy.NewProxy(r)
		p.Start(t)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
