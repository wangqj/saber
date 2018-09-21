package proxy

type Slot struct {
	id     int
	node   *Node
	status int
}

func NewSlot(id int, n *Node) *Slot {
	return &Slot{id, n, 1}
}

