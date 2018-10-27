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
}

type task struct {
	wg       *sync.WaitGroup
	Multi    []*redis.Resp
	response *redis.Resp
}

func NewSession(sock net.Conn) *Session {
	c := redis.NewConn(sock, 10000, 10000)
	s := Session{Conn: c}
	return &s
}

func (s *Session) Start(redisz *Router) {
	var ch = make(chan *task, 20480)
	go s.loopRead(ch)
	go s.loopWrite(ch)
}

func (s *Session) loopRead(ch chan *task) {
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

		r := &task{}
		r.Multi = multi
		r.wg = &sync.WaitGroup{}
		d := GetData()
		r.wg.Add(1)
		d.input <- r
		ch <- r
	}
}
func (s *Session) loopWrite(ch chan *task) {
	defer func() {
		log.Println("loopWrite close")
		s.Conn.Close()
	}()
	p := s.Conn.FlushEncoder()
	for i := 0; ; i++ {
		select {
		case t := <-ch:
			t.wg.Wait()
			p.Encode(t.response)
			p.Flush(true)

		}
	}
}
