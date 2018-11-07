package registry

import (
	"testing"
	"saber/proxy"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"golang.org/x/net/context"
	"time"
)


func TestEtcdx_LoadNodes(t *testing.T) {
	//o := utils.GetConf()
	ex := NewRegistryByPath("../config.toml")
	defer ex.Close()
	r := &proxy.Router{}
	ex.LoadNodes(r)
	for _, n := range r.Nodes {
		fmt.Println("TestEtcdx_LoadNodes ", n.Addr)
	}

}

func TestEtcdx_AddNode(t *testing.T) {
	ex := NewRegistryByPath("../config.toml")

	defer ex.Close()
	n := &proxy.Node{
		ID:        "3",
		Addr:      "127.0.0.1:6383",
		Status:    0,
		MaxIdle:   9,
		MaxActive: 11,
	}
	//n.BuildConn()
	ex.AddNode(n)
	//ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	//r, e := ex.CLi.Get(ctx, "/saber/nodes/"+n.ID)
	//fmt.Println(r, e)
	//b, _ := json.Marshal(n)
	//assert.JSONEq(t, "{\"ID\":\"1\",\"Addr\":\"127.0.0.1:6379\",\"Status\":1,\"MaxIdle\":0,\"MaxActive\":0}", string(b))
}

func TestEtcdx_LoadSlots(t *testing.T) {
	ex := NewRegistryByPath("../config.toml")
	defer ex.Close()
	r := &proxy.Router{}
	ex.LoadSlots(r)
}

func TestEtcdx_InitSlots(t *testing.T) {
	ex := NewRegistryByPath("../config.toml")
	defer ex.Close()
	r := &proxy.Router{}
	ex.LoadNodes(r)
	ex.InitSlots(r)
}

func TestEtcdx_ClearSlots(t *testing.T) {
	ex := NewRegistryByPath("../config.toml")
	defer ex.Close()
	ex.ClearSlots()
}

func TestEtcdx_WatchNodes(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	rch := cli.Watch(context.Background(), "/saber/nodes/3", clientv3.WithPrefix())
	wresp := <-rch
	fmt.Printf("wresp.Header.Revision: %d\n", wresp.Header.Revision)
	fmt.Println("wresp.IsProgressNotify:", wresp.IsProgressNotify())
	// wresp.Header.Revision: 0
	// wresp.IsProgressNotify: true
}

func TestEtcdx_RegProxy(t *testing.T) {
	ex := NewRegistryByPath("../config.toml")
	defer ex.Close()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		resp, _ := cli.Get(context.Background(), "/saber/proxy/"+proxy.GetPID(), clientv3.WithPrefix())

		for _, ev := range resp.Kvs {
			{
				fmt.Println(ev.Value)
			}
		}
	}
}