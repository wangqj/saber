package proxy

import (
	log "github.com/sirupsen/logrus"
	"net"
)

type Node struct {
	ID        string `json:"ID"`
	Addr      string `json:"Addr"`
	conn      net.Conn
	Status    int    `json:"Status"`
	MaxIdle   int    `json:"MaxIdle"`
	MaxActive int    `json:"MaxActive"`
}

type NodeStore struct {
	ID        string
	Addr      string
	Status    int
	MaxIdle   int
	MaxActive int
}

func (n *Node) BuildConn() (error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", n.Addr)
	if err != nil {
		log.Errorln("ResolveTCPAddr ", err)
		return err
	}
	n.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		n.Status = 0
		log.Errorln("DialTCP ", err)
		return err
	} else {
		n.Status = 1
	}
	if n.MaxIdle == 0 {
		n.MaxIdle = 10
	}
	if n.MaxActive == 0 {
		n.MaxActive = 5
	}
	log.Println("add redis node success ", n.Addr)
	return nil
}

/**
	根据。。。生成ID TODO
 */
func generateID() string {
	return "1"
}