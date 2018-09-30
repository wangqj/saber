package proxy

const OFFLINE int = 0 //?需要吗？
const ONLINE int = 1
const MIGRATE int = 2


type Slot struct {
	id     int
	node   *Node
	status int
}

func NewSlot(id int, n *Node) *Slot {
	return &Slot{id, n, 1}
}

func FindSlot() {

}
