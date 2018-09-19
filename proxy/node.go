package proxy

import (
	log "github.com/sirupsen/logrus"
	"net"
)

type Node struct {
	addr      string
	conn      net.Conn
	status    int
	MaxIdle   int
	MaxActive int
}

func (n *Node) newNode(addr string) (err error) {
	n.addr = addr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}
	n.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		n.status = 0
		return err
	} else {
		n.status = 1
	}
	log.Println("add redis success ", addr)
	return nil
}
