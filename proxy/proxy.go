package proxy

import (
	"net"
	"time"
	log "github.com/sirupsen/logrus"
	"os"
	"saber/utils"
	"strconv"
	"sync"
	"bytes"
)

type Proxy struct {
	redisz *Redisz
	status int
	addr   string
	exit struct {
		C chan struct{}
	}
	mu sync.Mutex
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
	p.mu.Lock()
	defer p.mu.Unlock()
	eh := make(chan error, 1)

	//建立socket，监听端口
	netListen, err := net.Listen("tcp", p.addr)
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer netListen.Close()
	t2 := time.Since(t)
	log.WithFields(log.Fields{"Spend time": t2, "listen tcp ": p.addr,}).Info("Successful startup! pid=", os.Getpid())

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
	proxyBuffer = bytes.TrimRight(proxyBuffer, "\x00")
	log.Println("request content is %s", string(proxyBuffer))

	//resp
	r2 := bytes.NewReader(proxyBuffer)

	data, err := ReadCommand(r2)
	if err != nil {
		log.Println(err)
	}
	log.Println("++++++++++++++++++++++++++")
	log.Println(data)
	log.Println(data.Args)
	log.Println(data.Name())
	log.Println(data.Value(0))
	log.Println(data.Value(1))
	log.Println(data.Format())
	log.Println("++++++++++++++++++++++++++")
	//slot
	s := redisz.GetSlot("")
	//
	////判断slot状态，
	//CheckSlot(s)
	//转发到redis
	conn, err := s.node.pool.Get()
	defer s.node.pool.Put(conn)
	//c := pool.GetConn()
	c := conn.(net.Conn)
	c.Write(proxyBuffer)
	//s.node.conn.Write(proxyBuffer)
	redisBuffer := make([]byte, 2048)
	re, readerr := c.Read(redisBuffer)
	if readerr != nil {
		log.Errorln("read redis error: %s", readerr.Error())
	}
	log.Println("get result ", string(redisBuffer[:re]), "---------------", "request content is %s", string(proxyBuffer))
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


