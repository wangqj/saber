package proxy

import (
	"saber/utils"
	log "github.com/sirupsen/logrus"
	"fmt"
)

type Router struct {
	Nodes []*Node
	Slots []*Slot
}

func (rz *Router) GetSlot(k string) *Slot {

	return rz.Slots[utils.HashCode(k)%1024]
}

func (rz *Router) GetNodeByNID(k string) *Node {

	for n, v := range rz.Nodes {
		fmt.Println(n, v)
		if v.ID == k {
			return v
		}
	}
	fmt.Println("mismatch", k)
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
