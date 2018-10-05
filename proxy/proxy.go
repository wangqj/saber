package proxy

import (
	"net"
	"time"
	log "github.com/sirupsen/logrus"
	"os"
	"saber/utils"
)

type Proxy struct {
	redisz *Redisz
	status int
	exit struct {
		C chan struct{}
	}
}

func NewProxy(o *utils.Option, r *Redisz) *Proxy {
	s := &Proxy{}
	s.redisz = r
	s.status = 1
	s.exit.C = make(chan struct{})
	return s
}

//启动proxy
func (p *Proxy) Start() {

	eh := make(chan error, 1)

	//建立socket，监听端口
	netListen, err := net.Listen("tcp", "localhost:16379")
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	//log.Println("listen port :" + strconv.Itoa(c.port))
	defer netListen.Close()
	log.Println("Waiting for clients")

	for {
		conn, err := netListen.Accept()
		if err != nil {
			time.Sleep(time.Duration(100 * time.Millisecond))
			log.Println(err)
			continue
		}
		log.Println(conn.RemoteAddr().String(), " tcp connect success")
		defer func() {
			eh <- err
		}()
		go handle(conn, p.redisz)
	}
	select {
	case <-p.exit.C:
		log.Warnf("[%p] proxy shutdown", p)
	case err := <-eh:
		log.Errorln(err, "[%p] proxy exit on error", p)
	}
}

func handle(proxyConn net.Conn, redisz *Redisz) {
	log.Println("handle request")
	proxyBuffer := make([]byte, 2048)
	proxyConn.Read(proxyBuffer)
	log.Println("request content is %s", proxyBuffer)
	//resp

	//slot
	s := redisz.GetSlot("")

	//判断slot状态，
	CheckSlot(s)
	//转发到redis
	s.node.conn.Write(proxyBuffer)
	redisBuffer := make([]byte, 2048)
	re, readerr := s.node.conn.Read(redisBuffer)
	if readerr != nil {
		log.Errorln("read redis error: %s", readerr.Error())
	}
	log.Println("get result ", string(redisBuffer[:re]))
	//返回结果
	proxyConn.Write(redisBuffer)
}

//TODO
func CheckSlot(s *Slot) {

	switch {
	case s.status == MIGRATE:
		log.Println("MIGRATE")
	case s.status == OFFLINE:
		log.Println("OFFLINE")
	case s.status == ONLINE:
		log.Println("ONLINE")

	}

}


