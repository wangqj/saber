package proxy

import (
	"saber/utils"
	log "github.com/sirupsen/logrus"
)

type Router struct {
	Nodes []*Node
	Slots []*Slot
}

func (rz *Router) GetSlot(k []byte) *Slot {
	return rz.Slots[utils.HashCode(k)%1024]
}

func (rz *Router) GetNodeByNID(k string) *Node {

	for _, v := range rz.Nodes {
		if v.ID == k {
			return v
		}
	}
	return nil
}

func (rz *Router) Dispatch(redisz *Router, v string) {

}

//TODO
func (rz *Router) CheckSlot(s *Slot) {

	switch {
	case s.Status == MIGRATE:
		log.Println("MIGRATE")
	case s.Status == OFFLINE:
		log.Println("OFFLINE")
	case s.Status == ONLINE:
		log.Println("ONLINE")

	}

}
