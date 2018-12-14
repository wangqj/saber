package app

import (
	"github.com/spf13/cobra"
	"time"
	"saber/utils"
	"saber/registry"
	_ "net/http/pprof"
	log "github.com/sirupsen/logrus"
	"saber/proxy"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var configFile string

var rootCmd = &cobra.Command{
	Use: "saber",
	Run: func(c *cobra.Command, args []string) {
		// 显示配置文件的名称和设置详细信息
		log.Printf("configFile: %s\nconfig: %#v", configFile, utils.Cfg)

		t := time.Now()
		log.Info("ready to start!")

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
	Version: "0.2",
}

func init() {
	//cobra.Command 执行前定义初始化处理。
	//rootCmd.Execute > 由于命令行参数处理按cobra.OnInitialize> rootCmd.Run的顺序执行，
	//使用标志接收的设置文件名读取配置文件，并在执行命令时使用设置文件的内容。
	cobra.OnInitialize(func() {
		// 将配置文件名定义为viper
		viper.SetConfigFile(configFile)
		// 读配置文件
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// 将配置文件的内容复制到结构中
		if err := viper.Unmarshal(&utils.Cfg); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})

	// 接收带有标志的设置文件名
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.toml", "config file name")
	rootCmd.MarkPersistentFlagRequired("config")
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
