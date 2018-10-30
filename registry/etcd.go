package registry

import (
	log "github.com/sirupsen/logrus"
	"saber/proxy"
	"github.com/coreos/etcd/clientv3"
	"time"
	"context"
	"saber/utils"
	"encoding/json"
	"strconv"
)

const SLOT_COUNT int = 16

type Etcdx struct {
	CLi *clientv3.Client

}

func NewRegistry() Registry {
	o := utils.GetConf()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{o.RegistryAdrr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Println("connect failed, err:", err)
		return nil
	}

	log.Println("connect etcd succuess ")
	//defer CLi.Close()
	r := &Etcdx{cli}
	return r
}

func NewRegistryByPath(path string) Registry {
	o := utils.GetConfByPath(path)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{o.RegistryAdrr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Println("connect failed, err:", err)
		return nil
	}

	log.Println("connect etcd succuess ")
	//defer CLi.Close()
	r := &Etcdx{cli}
	return r
}

func (e *Etcdx) LoadNodes(r *proxy.Router) {
	//设置1秒超时，访问etcd有超时控制
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	resp, err := e.CLi.Get(ctx, "/saber/nodes/", clientv3.WithPrefix())
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
			r.Nodes = append(r.Nodes, &p)
		} else {
			log.Println("this node status is :", p.Status)
		}

	}
}

func (e *Etcdx) LoadSlots(r *proxy.Router) {
	//设置1秒超时，访问etcd有超时控制
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	resp, err := e.CLi.Get(ctx, "/saber/slots/", clientv3.WithPrefix())
	if err != nil {
		log.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		log.Printf("%s : %s\n", ev.Key, ev.Value)
		var s proxy.Slot
		err := json.Unmarshal(ev.Value, &s)

		if err != nil {
			log.Println("json unmarshal failed, err:", err)
			continue
		}
		if s.Status == 1 {
			s.Node = r.GetNodeByNID(s.NID)
			r.Slots = append(r.Slots, &s)
		} else {
			log.Println("this slot status is :", s.Status)
		}
	}
}

func (e *Etcdx) Close() {
	e.CLi.Close()
}

func (e *Etcdx) AddNode(n *proxy.Node) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	b, err := json.Marshal(n)
	if err != nil {
		log.Println(err)
	} else {
		resp, ce := e.CLi.Put(ctx, "/saber/nodes/"+n.ID, string(b))
		if ce != nil {
			log.Println(ce)
		} else {
			log.Println(resp)
		}
	}
}

func (e *Etcdx) InitSlots(r *proxy.Router) {
	if r.Nodes == nil || len(r.Nodes) == 0 {
		log.Errorln("nodes is nil")
		return
	}
	if len(r.Slots) == 0 {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		for i := 1; i <= 1024; i++ {
			s := proxy.NewSlot(strconv.Itoa(i), r.Nodes[i%len(r.Nodes)])
			b, err := json.Marshal(s)
			if err != nil {
				log.Println(err)
			} else {
				log.Info("init slot ", *s)
				resp, ce := e.CLi.Put(ctx, "/saber/slots/"+string(s.ID), string(b))
				if ce != nil {
					log.Println(ce)
				} else {
					log.Println(resp)
				}
			}

		}
	}
}

func (e *Etcdx) ClearSlots() {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)

	resp, ce := e.CLi.Delete(ctx, "/saber/slots/", clientv3.WithPrefix())

	if ce != nil {
		log.Println(ce)
	} else {
		log.Println(resp)
	}

}

func RegProxy() {

}