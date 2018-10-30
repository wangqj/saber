package registry

import (
	"testing"
	"saber/proxy"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func newEtcdx() *Etcdx {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return nil
	}

	fmt.Println("connect etcd succuess ")
	//defer CLi.Close()
	r := &Etcdx{cli}

	return r
}

func TestEtcdx_LoadNodes(t *testing.T) {
	//o := utils.GetConf()
	ex := NewRegistry()
	defer ex.Close()
	r := proxy.Router{}
	ex.LoadNodes(&r)
	for _, n := range r.Nodes {
		fmt.Println("TestEtcdx_LoadNodes ", n.Addr)
	}

}

func TestEtcdx_AddNode(t *testing.T) {
	ex := newEtcdx()

	defer ex.Close()
	n := &proxy.Node{
		ID:        "2",
		Addr:      "127.0.0.1:6382",
		Status:    1,
		MaxIdle:   10,
		MaxActive: 3,
	}
	//n.BuildConn()
	ex.AddNode(n)
	//ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	//r, e := ex.CLi.Get(ctx, "/saber/nodes/"+n.ID)
	//fmt.Println(r, e)
	//b, _ := json.Marshal(n)
	//assert.JSONEq(t, "{\"ID\":\"1\",\"Addr\":\"127.0.0.1:6379\",\"Status\":1,\"MaxIdle\":0,\"MaxActive\":0}", string(b))
}


