package registry

import (
	"saber/proxy"
)

type Registry interface {
	//NewRegistry(o *utils.Option) (*Registry)
	LoadNodes(r *proxy.Router) //(error)
	LoadSlots(r *proxy.Router) //(error)
	Close()
}
