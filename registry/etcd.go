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

type Etcdx struct {
	cli *clientv3.Client
}

func NewEtcdx(o *utils.Option) *Etcdx {
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

func (e *Etcdx) LoadNodes(r *proxy.Redisz) {
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
		n := proxy.Node{}
		n.NewNode(string(ev.Value))
		r.Nodes = append(r.Nodes, &n)
	}
}

func (e *Etcdx) LoadSlots(r *proxy.Redisz) {
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

//func exist(cli *clientv3.Client,key string)  {
//	log.Println("获取值")
//	if resp, err := cli.Get(context.TODO(), key); err != nil {
//		log.Fatal(err)
//	} else {
//		log.Println("resp: ", resp)
//		resp.Kvs
//	}
//}

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