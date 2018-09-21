package proxy

import (
	"net"
	"rocket/redisy"
	"time"
	log "github.com/sirupsen/logrus"
	"saber/utils"
	"os"
)

type Proxy struct {
	exit struct {
		C chan struct{}
	}
}

func NewProxy() *Proxy {

	return nil
}

//启动proxy
func (p *Proxy) Start(o *utils.Option) {

	eh := make(chan error, 1)

	//建立socket，监听端口
	netListen, err := net.Listen("tcp", "localhost:6379")
	CheckError(err)
	//log.Println("listen port :" + strconv.Itoa(c.port))
	defer netListen.Close()
	log.Println("Waiting for clients")
	rz := redisy.Redisz{}
	rz.Init()
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
