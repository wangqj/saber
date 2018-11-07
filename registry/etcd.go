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
	"strings"
	"fmt"
)

const SLOT_COUNT int = 16

type Etcdx struct {
	CLi *clientv3.Client

}

func NewRegistry() Registry {
	o := utils.GetConf()
	var addrs []string
	if strings.Index(o.RegistryAdrr, ",") > 0 {
		addrs = strings.Split(o.RegistryAdrr, ",")
	} else {
		addrs = append(addrs, o.RegistryAdrr)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
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
	var addrs []string
	if strings.Index(o.RegistryAdrr, ",") > 0 {
		addrs = strings.Split(o.RegistryAdrr, ",")
	} else {
		addrs = append(addrs, o.RegistryAdrr)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
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
		log.Errorln("get nodes failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		log.Printf("load Nodes %s : %s\n", ev.Key, ev.Value)
		var p proxy.Node
		err := json.Unmarshal(ev.Value, &p)
		if err != nil {
			log.Errorln("json unmarshal failed, err:", err)
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
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	resp, err := e.CLi.Get(ctx, "/saber/slots/", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		log.Errorln("get slots failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		//log.Printf("%s : %s\n", ev.Key, ev.Value)
		var s proxy.Slot
		err := json.Unmarshal(ev.Value, &s)

		if err != nil {
			log.Errorln("json unmarshal failed, err:", err)
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
				resp, ce := e.CLi.Put(ctx, "/saber/slots/"+utils.AddZeroForStr(s.ID, 4), string(b))
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

func (e *Etcdx) WatchNodes(r *proxy.Router) {
	rch := e.CLi.Watch(context.Background(), "/saber/nodes/", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				var p proxy.Node
				err := json.Unmarshal(ev.Kv.Value, &p)
				if err != nil {
					log.Errorln("json unmarshal failed, err:", err)
				}
				log.Println("Addr=", p.Addr)
				//这里很复杂，分多种情况
				f := false
				for k, v := range r.Nodes {
					if v.ID == p.ID {
						if p.Status == 0 {
							//没有处理slot，如果是node挂掉了，应该用备份节点替换，如果是手动操作，应该先迁移slot
							//此处应该认为此node没有slot了
							r.Nodes[k].Close()
							r.Nodes = append(r.Nodes[:k], r.Nodes[k+1:]...)
							//v=&p
							f = true
							break
						}
					}
				}
				if !f {
					if p.Status == 1 {
						p.BuildConn()
						r.Nodes = append(r.Nodes, &p)
					} else {
						//TODO
						log.Println("this node status is :", p.Status)
					}
				}

			case clientv3.EventTypeDelete:
				log.Println("EventTypeDelete :")
			default:
				log.Println("defult :", ev.Kv)

			}

		}
	}
}

func (e *Etcdx) WatchSlots(r *proxy.Router) {
	rch := e.CLi.Watch(context.Background(), "/saber/slots/", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				var slot proxy.Slot
				err := json.Unmarshal(ev.Kv.Value, &slot)
				if err != nil {
					log.Errorln("json unmarshal failed, err:", err)
				}
				log.Println("Addr=", slot)
				for _, v := range r.Slots {
					if slot.ID == v.ID {
						v = &slot
					}
				}
			case clientv3.EventTypeDelete:
				log.Println("EventTypeDelete :")
			default:
				log.Println("defult :", ev.Kv)

			}

		}
	}
}

func (e *Etcdx) LoadAll(r *proxy.Router) {
	e.LoadNodes(r)
	e.LoadSlots(r)
	go e.WatchNodes(r)
	go e.WatchSlots(r)
	e.RegProxy()
}

//注册proxy
func (e *Etcdx) RegProxy() {
	grantResp, err := e.CLi.Grant(context.TODO(), 10)
	if err != nil {
		log.Errorln(err)
		return
	}
	keepResp, err := e.CLi.Put(context.TODO(), "/saber/proxy/"+proxy.GetPID(), "living", clientv3.WithLease(grantResp.ID))
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Println(keepResp)
	go keepAlive(e, grantResp)
}

//每隔5秒心跳一次
func keepAlive(e *Etcdx, grantResp *clientv3.LeaseGrantResponse) {
	for {
		time.Sleep(time.Duration(5) * time.Second)
		_, err := e.CLi.KeepAliveOnce(context.TODO(), grantResp.ID)
		if err != nil {
			log.Errorln(err)
		}
	}
}