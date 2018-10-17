package proxy

import (
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/gomodule/redigo/redis"
)

type Node struct {
	ID        string `json:"ID"`
	Addr      string `json:"Addr"`
	//conn      net.Conn
	pool      *redis.Pool
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

	p := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", n.Addr)
			if err != nil {
				log.Println("00000000", err)
				return nil, err
			}
			redis.DialConnectTimeout(3 * time.Second)
			return c, nil
		},
		MaxActive:       100,
		MaxIdle:         100,
		Wait:            true,
		MaxConnLifetime: 60 * time.Second,
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
