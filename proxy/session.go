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
	go s.rhandle(&ch, redisz)
	go s.whandle(&ch)
}

func (s *Session) rhandle(ch *chan task, redisz *Router) {
	defer func() {
		fmt.Println("rhandle close")
	}()
	for {
		s.conn.SetReadDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
		proxyBuffer := make([]byte, 2048)
		n, err := s.conn.Read(proxyBuffer)
		if err != nil {
			fmt.Println(os.Stderr, "3Fatal error: %s", err.Error())
			s.conn.Close()
			break
			//os.Exit(1)
		}
		r2 := bytes.NewReader(proxyBuffer[:n])

		data, err := ReadCommand(r2)
		input := data.Value(1)
		//fmt.Println("rhandle =",input)
		t := task{req: input}
		*ch <- t
	}
}
func (s *Session) whandle(ch *chan task) {
	defer func() {
		fmt.Println("whandle close")
		s.conn.Close()
	}()
	for i := 0; ; i++ {
		select {
		case msg := <-*ch:
			//fmt.Println("msg=",msg)
			d := "+" + msg.req + "\r\n"
			//fmt.Println(d)
			_, err := s.conn.Write([]byte(d))
			if err != nil {
				fmt.Println(os.Stderr, "whandle error: %s", err.Error())
				break
			}
		}
	}
}
