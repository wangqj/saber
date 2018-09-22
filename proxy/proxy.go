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
	netListen, err := net.Listen("tcp", "localhost:6379")
	CheckError(err)
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
		go handle(conn)
	}
	select {
	case <-p.exit.C:
		log.Warnf("[%p] proxy shutdown", p)
	case err := <-eh:
		log.Errorln(err, "[%p] proxy exit on error", p)
	}
}

func handle(conn net.Conn) {
	//resp

	//slot

	//判断slot状态，做处理

	//转发到redis

	//返回结果
}

func CheckError(err error) {
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
