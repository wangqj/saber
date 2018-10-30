package proxy

const OFFLINE int = 0 //?需要吗？
const ONLINE int = 1
const MIGRATE int = 2


type Slot struct {
	ID     string `json:"ID"`
	Node   *Node
	NID    string `json:"NID"`
	Status int    `json:"Status"`
}

func NewSlot(id string, n *Node) *Slot {
	return &Slot{id, n, n.ID, 1}
}

func FindSlot() {

}
