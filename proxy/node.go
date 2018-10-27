package proxy

import (
	log "github.com/sirupsen/logrus"
)

type Node struct {
	ID        string `json:"ID"`
	Addr      string `json:"Addr"`
	Status    int    `json:"Status"`
	MaxIdle   int    `json:"MaxIdle"`
	MaxActive int    `json:"MaxActive"`
}

func init() {
	log.Println("init run!")
}


/**
	根据。。。生成ID TODO
 */
func generateID() string {
	return "1"
}
