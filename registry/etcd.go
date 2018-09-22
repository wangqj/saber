package registry

import (
	log "github.com/sirupsen/logrus"
	"saber/proxy"
	"github.com/coreos/etcd/clientv3"
	"time"
	"context"
	"saber/utils"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := e.cli.Get(ctx, "/saber/nodes/")
	cancel()
	if err != nil {
		log.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		log.Printf("%s : %s\n", ev.Key, ev.Value)
		n := proxy.Node{}
		n.NewNode(string(ev.Value))
		r.Nodes = append(r.Nodes, &n)
	}
}

func (e *Etcdx) LoadSlots(r *proxy.Redisz) {
	//设置1秒超时，访问etcd有超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := e.cli.Get(ctx, "/saber/slots/")
	cancel()
	if err != nil {
		log.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		log.Printf("%s : %s\n", ev.Key, ev.Value)
		//TODO
	}
}
