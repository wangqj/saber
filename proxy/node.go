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
	conn      *redis.Conn
	proc      *Processor
}

func (n *Node) BuildConn() (error) {
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
	n.proc = d
}

// TODO
func (n *Node) keepAlive() {

}
