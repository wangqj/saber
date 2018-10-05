package proxy

import (
	"testing"
	"github.com/sirupsen/logrus"
)

func TestNode_BuildConn(t *testing.T) {
	n := Node{
		ID:        "11",
		Addr:      "127.0.0.1:6379",
		Status:    1,
		MaxIdle:   10,
		MaxActive: 3,
	}
	n.BuildConn()
	logrus.Print(n)
}
