package proxy

import (
	"time"
	"saber/proxy/redis"
	log "github.com/sirupsen/logrus"
)

type Node struct {
	ID        string `json:"ID"`
	Addr      string `json:"Addr"`
	Status    int    `json:"Status"`
	MaxIdle   int    `json:"MaxIdle"`
	MaxActive int    `json:"MaxActive"`
	SlotCount int    `json:"SlotCount"`
	conn      *redis.Conn
	processor *Processor
}

//创建连接，每个node一个连接
func (n *Node) BuildConn() (error) {
	//TODO buffer可配置
	c, err := redis.DialTimeout(n.Addr, time.Second*5,
		128*1024,
		128*1024)
	if err != nil {
		return err
	}
	n.conn = c
	log.Info("BuildConn success. ", c.Sock)
	n.buildProcessor()
	go n.keepAlive()
	return nil
}

func (n *Node) buildProcessor() {
	d := NewProcessor(n.conn)
	go d.Start()
	n.processor = d
}

// TODO
func (n *Node) keepAlive() {

}

func (n *Node) Close() {
	//TODO 此处应该检查是否还有属于此Node的slot
	n.processor.Stop()
	n.conn.Close()
}
