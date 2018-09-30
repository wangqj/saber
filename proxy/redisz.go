package proxy

import "saber/utils"

type Redisz struct {
	Nodes []*Node
	Slots []*Slot
}

func (rz *Redisz) GetSlot(k string) *Slot {
	return rz.Slots[utils.HashCode(k)%1024]
}
