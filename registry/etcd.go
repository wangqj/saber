package registry

import (
	log "github.com/sirupsen/logrus"
	"saber/proxy"
	"github.com/coreos/etcd/clientv3"
	"time"
	"context"
	"saber/utils"
	"encoding/json"
)

const SLOT_COUNT int = 16

type Etcdx struct {
	cli *clientv3.Client

}

func NewRegistry(o *utils.Option) Registry {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{o.RegistryAdrr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Println("connect failed, err:", err)
		return nil
	}

	log.Println("connect etcd succuess ")
	//defer cli.Close()
	r := &Etcdx{cli}
	return r
}

func (e *Etcdx) LoadNodes(r *proxy.Router) {
	//设置1秒超时，访问etcd有超时控制
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	resp, err := e.cli.Get(ctx, "/saber/nodes/", clientv3.WithPrefix())
	if err != nil {
		log.Errorln(err)
		log.Fatal(err)
	} else {
		log.Println("resp: ", resp)
	}
	//cancel()
	if err != nil {
		log.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		log.Printf("load Nodes %s : %s\n", ev.Key, ev.Value)
		var p proxy.Node
		err := json.Unmarshal(ev.Value, &p)
		if err != nil {
			log.Println("json unmarshal failed, err:", err)
		}
		log.Println("Addr=", p.Addr)
		if p.Status == 1 {
			p.BuildConn()
			//n, err := proxy.NewNode(string(p.Addr))

			if err != nil {
				log.Println("init node failed, err:", err)
			} else {
				r.Nodes = append(r.Nodes, &p)
			}
		} else {
			log.Println("this node status is :", p.Status)
		}

	}
}

func (e *Etcdx) LoadSlots(r *proxy.Router) {
	//TODO
	for i := 0; i < 1024; i++ {
		s := proxy.NewSlot(i, r.Nodes[i%len(r.Nodes)])
		r.Slots = append(r.Slots, s)
	}

	//设置1秒超时，访问etcd有超时控制
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//resp, err := e.cli.Get(ctx, "/saber/slots/")
	//cancel()
	//if err != nil {
	//	log.Println("get failed, err:", err)
	//	return
	//}
	//for _, ev := range resp.Kvs {
	//	log.Printf("%s : %s\n", ev.Key, ev.Value)
	//	//TODO
	//}
}

func (e *Etcdx) Close() {
	e.cli.Close()
}

func (e *Etcdx) AddNode(n proxy.Node) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	b, err := json.Marshal(n)
	if err != nil {
		log.Println(err)
	} else {
		resp, ce := e.cli.Put(ctx, "/saber/nodes/"+n.ID, string(b))
		if ce != nil {
			log.Println(ce)
		} else {
			log.Println(resp)
		}
	}
}

func RegProxy() {

}