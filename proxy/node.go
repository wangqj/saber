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

func NewNode(addr string) (*Node, error) {
	n := Node{}
	n.ID = generateID()
	n.Addr = addr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Errorln("ResolveTCPAddr ", err)
		return nil, err
	}
	n.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		n.Status = 0
		log.Errorln("DialTCP ", err)
		return nil, err
	} else {
		n.Status = 1
	}
	//TODO,先写死默认值
	n.MaxIdle = 10
	n.MaxActive = 5
	log.Println("add redis node success ", addr)
	return &n, nil
}

/**
	根据。。。生成ID TODO
 */
func generateID() string {
	return "1"
}