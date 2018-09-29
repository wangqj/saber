package proxy

import "saber/utils"

type Redisz struct {
	Nodes []*Node
	Slots []*Slot
}

func (rz *Redisz) get(k string) interface{} {
	return rz.Slots[utils.HashCode(k)%1024]
}
