package proxy

import (
	log "github.com/sirupsen/logrus"
	"net"
	"time"
	"saber/utils"
)

type Node struct {
	ID        string `json:"ID"`
	Addr      string `json:"Addr"`
	//conn      net.Conn
	pool      utils.Pool
	Status    int    `json:"Status"`
	MaxIdle   int    `json:"MaxIdle"`
	MaxActive int    `json:"MaxActive"`
}

func init() {
	log.Println("init run!")
}

func (n *Node) BuildConn() (error) {
	if n.MaxIdle == 0 {
		n.MaxIdle = 10
	}
	if n.MaxActive == 0 {
		n.MaxActive = 5
	}

	//factory 创建连接的方法
	factory := func() (interface{}, error) { return net.Dial("tcp", n.Addr) }

	//close 关闭连接的方法
	close := func(v interface{}) error { return v.(net.Conn).Close() }

	//ping 检测连接的方法
	//ping := func(v interface{}) error { return nil }

	//创建一个连接池： 初始化5，最大连接30
	poolConfig := &utils.PoolConfig{
		Name:       n.Addr,
		InitialCap: 5,
		MaxCap:     30,
		Factory:    factory,
		Close:      close,
		//Ping:       ping,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}
	p, err := utils.NewChannelPool(poolConfig)
	if err != nil {
		log.Errorln("NewGenericPool ", err)
		return err
	}
	n.pool = p
	return nil
}

/**
	根据。。。生成ID TODO
 */
func generateID() string {
	return "1"
}
