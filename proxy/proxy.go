package proxy

import (
	"net"
	"time"
	log "github.com/sirupsen/logrus"
	"os"
	"saber/utils"
	"strconv"
)

type Proxy struct {
	redisz *Redisz
	status int
	addr   string
	exit struct {
		C chan struct{}
	}
}

func NewProxy(o *utils.Option, r *Redisz) *Proxy {
	s := &Proxy{}
	s.redisz = r
	s.status = 1
	s.addr = "127.0.0.1:" + strconv.Itoa(o.Port)
	s.exit.C = make(chan struct{})
	return s
}

//启动proxy
func (p *Proxy) Start(t time.Time) {

	eh := make(chan error, 1)

	//建立socket，监听端口
	netListen, err := net.Listen("tcp", p.addr)
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer netListen.Close()
	t2 := time.Since(t)
	log.WithFields(log.Fields{"Spend time": t2, "listen tcp ": p.addr,}).Info("Successful startup!")

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


