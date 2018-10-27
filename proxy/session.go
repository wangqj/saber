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

func (s *Session) Start(redisz *Router) {
	log.Println("session start")
	go s.loopRead()
	go s.loopWrite()
}

func (s *Session) loopRead() {
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

		r := &task{
			Multi: multi,
			wg:    &sync.WaitGroup{},
		}
		d := GetData()
		r.wg.Add(1)
		d.input <- r
		s.ch <- r
	}
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
