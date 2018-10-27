package registry

import (
	"testing"
	"saber/utils"
	"saber/proxy"
	"time"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

type gloable struct {
	etcdx Etcdx
}

func TestNewEtcdxNewEtcdx(t *testing.T) {

}

func TestEtcdx_LoadNodes(t *testing.T) {
	o := utils.GetConf()
	ex := NewEtcdx(o)
	defer ex.Close()
	r := proxy.Router{}
	ex.LoadNodes(&r)
	for _, n := range r.Nodes {
		fmt.Println("TestEtcdx_LoadNodes ", n.Addr)
	}

}

func TestEtcdx_AddNode(t *testing.T) {
	o := utils.GetConf()
	ex := NewEtcdx(o)
	defer ex.Close()
	n := proxy.Node{
		ID:        "1",
		Addr:      "127.0.0.1:6381",
		Status:    1,
		MaxIdle:   10,
		MaxActive: 3,
	}
	n.BuildConn()
	ex.AddNode(n)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	r, e := ex.cli.Get(ctx, "/saber/nodes/"+n.ID)
	fmt.Println(r, e)
	b, _ := json.Marshal(n)
	assert.JSONEq(t, "{\"ID\":\"1\",\"Addr\":\"127.0.0.1:6379\",\"Status\":1,\"MaxIdle\":0,\"MaxActive\":0}", string(b))
}
