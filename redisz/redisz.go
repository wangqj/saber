package redisz

import (
	"saber/proxy"
)

type Redisz struct {
	Nodes []*proxy.Node
	Slots []*proxy.Slot
}

func (ns *Redisz) get(k string) interface{} {

	return nil
}
