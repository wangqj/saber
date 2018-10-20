package proxy

import (
	"saber/utils"
	log "github.com/sirupsen/logrus"
	"bytes"
	"net"
)

type Router struct {
	Nodes []*Node
	Slots []*Slot
}

func (rz *Router) GetSlot(k string) *Slot {
	return rz.Slots[utils.HashCode(k)%1024]
}
func (rz *Router) Dispatch(redisz *Router, v string) {

}

func (rz *Router) HandleRequest(ch *chan task, task *task) {
	reader := bytes.NewReader(task.reqeust.data)
	data, err := ReadCommand(reader)
	if err != nil {
		log.Errorln("HandleRequest resp error: %s", err.Error())
	}
	//slot
	s := rz.GetSlot("")

	conn, err := s.node.pool.Get()
	if err != nil {
		log.Errorln("get pool error: %s", err.Error())
	}

	c := conn.(net.Conn)

	switch {
	case data.Name() == "get":
		c.Write(task.reqeust.data)
	case data.Name() == "set":
		c.Write(task.reqeust.data)
	default:
		log.Println("no case ,", data.Name())

	}
	go write(ch, task, c, s)
}

func write(ch *chan task, task *task, c net.Conn, s *Slot) {
	defer s.node.pool.Put(c)

	redisBuffer := make([]byte, 2048)
	n, readerr := c.Read(redisBuffer)

	if readerr != nil {
		log.Errorln("read redis error: %s", readerr.Error())
	}

	rr := bytes.NewReader(redisBuffer[:n])
	cmd, _ := ReadCommand(rr)
	resp := &Resp{
		data: redisBuffer[:n],
		cmd:  cmd,
	}
	task.response = resp
	*ch <- *task
}

func loopReadRedis(task *task) {

}

//TODO
func (rz *Router) CheckSlot(s *Slot) {

	switch {
	case s.status == MIGRATE:
		log.Println("MIGRATE")
	case s.status == OFFLINE:
		log.Println("OFFLINE")
	case s.status == ONLINE:
		log.Println("ONLINE")

	}

}
