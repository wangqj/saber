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
