package proxy

import (
	"saber/utils"
	log "github.com/sirupsen/logrus"
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

func (rz *Router) Handle() {
	//slot
	s := rz.GetSlot("")
	////判断slot状态，
	//CheckSlot(s)
	//转发到redis
	conn := s.node.pool.Get()
	defer conn.Close()
	//conn.Do()
	//c.Write(proxyBuffer)
	//redisBuffer := make([]byte, 2048)
	//re, readerr := c.Read(redisBuffer)
	//if readerr != nil {
	//	log.Errorln("read redis error: %s", readerr.Error())
	//}

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
