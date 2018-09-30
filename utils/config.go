package utils

type Option struct {
	Port         int
	NCPU         int
	RegistryAdrr string
}

//读取本地配置文件TODO

func LoadConf() (*Option) {
	o := &Option{}
	o.RegistryAdrr = "127.0.0.1:2379"
	return o
}
