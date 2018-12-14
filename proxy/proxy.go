package proxy

import (
	"net"
	"time"
	log "github.com/sirupsen/logrus"
	"os"
	"saber/utils"
	"strconv"
	"sync"
)

type Proxy struct {
	router *Router
	status int
	addr   string
	exit struct {
		C chan struct{}
	}
	mu sync.Mutex
}

func NewProxy(r *Router) *Proxy {
	o := utils.Cfg
	s := &Proxy{}
	s.router = r
	s.status = 1
	s.addr = "0.0.0.0:" + strconv.Itoa(o.Port)
	s.exit.C = make(chan struct{})
	return s
}

//启动proxy one loop per thread
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
			log.Println(err)
			continue
		}
		log.Println(conn.RemoteAddr().String(), " tcp connect success")
		defer func() {
			eh <- err
		}()
		s := NewSession(conn)
		s.Start(p.router)

	}
	select {
	case <-p.exit.C:
		log.Warnf("[%p] proxy shutdown", p)
	case err := <-eh:
		log.Errorln(err, "[%p] proxy exit on error", p)
	}
}

func GetPID() string {
	o := utils.Cfg
	pid := utils.GetIntranetIp() + ":" + strconv.Itoa(o.Port)
	return pid
}




