package proxy

import (
	"os"
	"sync"
	"saber/proxy/redis"
	"net"
	log "github.com/sirupsen/logrus"
	"fmt"
)

type Session struct {
	Conn *redis.Conn
	ch   chan *task
}

type task struct {
	wg       *sync.WaitGroup
	Multi    []*redis.Resp
	response *redis.Resp
}

func NewSession(sock net.Conn, ) *Session {
	//TODO，rbuf需要从配置中取
	c := redis.NewConn(sock, 10000, 10000)
	s := &Session{Conn: c, ch: make(chan *task, 20480)}
	return s
}

func (s *Session) Start(router *Router) {
	log.Println("session start")
	go s.loopRead(router)
	go s.loopWrite()
}

//循环读取数据
func (s *Session) loopRead(router *Router) {
	defer func() {
		log.Println("loopRead close")
	}()
	for {
		multi, err := s.Conn.DecodeMultiBulk()
		if err != nil {
			log.Println(os.Stderr, "loopRead error: %s", err.Error())
			s.Conn.Close()
			break
		}
		r := buildTask(multi)
		handleRequest(r, router)
		s.ch <- r
	}
}

//构建任务数据
func buildTask(respArray []*redis.Resp) *task {
	r := &task{
		Multi: respArray,
		wg:    &sync.WaitGroup{},
	}
	r.wg.Add(1)
	return r
}

//处理请求，主要负责解析resp，分发请求
func handleRequest(t *task, router *Router) {

	switch getFlag(t.Multi) {
	case "ping":
		//TODO
		fmt.Println("ping")
	default:
		//根据muti获取slot
		resp := t.Multi[1]
		p := router.GetSlot(resp.Value).Node.processor
		//传递到processor通道
		p.input <- t
	}

}

//最后处理的结果返回给客户端
func (s *Session) loopWrite() {
	defer func() {
		log.Println("loopWrite close")
		s.Conn.Close()
	}()
	p := s.Conn.FlushEncoder()
	for {
		select {
		case t := <-s.ch:
			t.wg.Wait()
			p.Encode(t.response)
			p.Flush(true)
		}
	}
}
