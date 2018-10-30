package registry

import (
	"saber/proxy"
)


type Registry interface {

	LoadNodes(r *proxy.Router) //(error)
	LoadSlots(r *proxy.Router) //(error)
	AddNode(n *proxy.Node)
	InitSlots(r *proxy.Router)
	ClearSlots()
	Close()
}
