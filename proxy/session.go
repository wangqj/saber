package proxy

import (
	"fmt"
	"net"
	"time"
	"os"
	"bytes"
)

type Session struct {
	conn net.Conn
}

type Request struct {
}
type task struct {
	req string
	res string
}

func NewSession(c net.Conn) *Session {
	s := Session{conn: c}
	return &s
}

func (s *Session) Start(redisz *Router) {
	var ch = make(chan task)
	go s.loopRead(&ch, redisz)
	go s.loopWrite(&ch)
}

func (s *Session) loopRead(ch *chan task, redisz *Router) {
	defer func() {
		fmt.Println("loopRead close")
	}()
	for {
		s.conn.SetReadDeadline(time.Now().Add(time.Duration(100) * time.Second))
		proxyBuffer := make([]byte, 2048)
		n, err := s.conn.Read(proxyBuffer)
		if err != nil {
			fmt.Println(os.Stderr, "3Fatal error: %s", err.Error())
			s.conn.Close()
			break
		}
		r2 := bytes.NewReader(proxyBuffer[:n])

		data, err := ReadCommand(r2)
		input := data.Value(1)
		t := task{req: input}
		*ch <- t
	}
}
func (s *Session) loopWrite(ch *chan task) {
	defer func() {
		fmt.Println("loopWrite close")
		s.conn.Close()
	}()
	for i := 0; ; i++ {
		select {
		case msg := <-*ch:
			d := "+" + msg.req + "\r\n"
			_, err := s.conn.Write([]byte(d))
			if err != nil {
				fmt.Println(os.Stderr, "loopWrite error: %s", err.Error())
				break
			}
		}
	}
}
