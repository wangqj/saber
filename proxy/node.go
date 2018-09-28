package proxy

import (
	log "github.com/sirupsen/logrus"
	"net"
)

type Node struct {
	ID        string
	Addr      string
	conn      net.Conn
	Status    int
	MaxIdle   int
	MaxActive int
}

func (n *Node) NewNode(addr string) (err error) {
	n.ID = generateID()
	n.Addr = addr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}
	n.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		n.Status = 0
		return err
	} else {
		n.Status = 1
	}
	log.Println("add redis success ", addr)
	return nil
}

/**
	根据。。。生成ID TODO
 */
func generateID() string {
	return "1"
}