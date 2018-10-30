package proxy

import (
	"os"
	"sync"
	"saber/proxy/redis"
	"net"
	log "github.com/sirupsen/logrus"
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
	c := redis.NewConn(sock, 10000, 10000)
	s := Session{Conn: c, ch: make(chan *task, 20480)}
	return &s
}

func (s *Session) Start(router *Router) {
	log.Println("session start")
	go s.loopRead(router)
	go s.loopWrite()
}

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

func buildTask(respArray []*redis.Resp) *task {
	r := &task{
		Multi: respArray,
		wg:    &sync.WaitGroup{},
	}
	r.wg.Add(1)
	return r
}

func handleRequest(t *task, router *Router) {
	//TODO 根据muti获取slot
	p := router.GetSlot("").Node.proc
	p.input <- t
}

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
