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

type NodePool struct {
	Name     string
	Conn     net.Conn
	activeAt time.Time
}

func (p *NodePool) GetActiveTime() time.Time {
	return p.activeAt
}
func (p *NodePool) Close() error {
	log.Println(p.Name, " closed")
	p.Conn.Close()
	return nil
}
func (p *NodePool) GetConn() net.Conn {
	return p.Conn
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
	p, err := utils.NewGenericPool(1, n.MaxActive, time.Minute*10, func() (utils.Poolable, error) {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", n.Addr)
		if err != nil {
			log.Errorln("ResolveTCPAddr ", err)
			return nil, err
		}
		c, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			n.Status = 0
			log.Errorln("DialTCP ", err)
			return nil, err

		} else {
			n.Status = 1
		}
		log.Println("get redis conn success ", n.Addr)
		return &NodePool{Name: n.ID, Conn: c, activeAt: time.Now()}, nil
	})
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
