package proxy

type Node struct {
	ID        string `json:"ID"`
	Addr      string `json:"Addr"`
	Status    int    `json:"Status"`
	MaxIdle   int    `json:"MaxIdle"`
	MaxActive int    `json:"MaxActive"`
}
