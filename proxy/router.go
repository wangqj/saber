package proxy

import (
	"saber/utils"
	log "github.com/sirupsen/logrus"
	"time"
	"saber/proxy/redis"
	"fmt"
)

type Router struct {
	Nodes []*Node
	Slots []*Slot
}

func (rz *Router) GetSlot(k string) *Slot {
	return rz.Slots[utils.HashCode(k)%1024]
}
func (rz *Router) Dispatch(redisz *Router, v string) {

}

func (rz *Router) HandleRequest(ch chan *task, task *task) {
	fmt.Println("-----1--", task.Multi)
	c, err := redis.DialTimeout(":6381", time.Second*5,
		128,
		128)
	if err != nil {
		log.Errorln("HandleRequest 1resp error: %s", err.Error())
	}
	p := c.FlushEncoder()
	p.MaxInterval = time.Millisecond
	p.MaxBuffered = cap(ch) / 2
	err = p.EncodeMultiBulk(task.Multi)
	if err != nil {
		log.Errorln("HandleRequest 2resp error: %s", err.Error())
	}

	resp, err := c.Decode()
	if err != nil {
		log.Errorln("HandleRequest 3resp error: %s", err.Error())
	}
	task.response = resp

	//reader := bytes.NewReader(task.reqeust.data)
	//data, err := ReadCommand(reader)
	//if err != nil {
	//	log.Errorln("HandleRequest resp error: %s", err.Error())
	//}
	////slot
	//s := rz.GetSlot("")
	//
	//conn, err := s.node.pool.Get()
	//if err != nil {
	//	log.Errorln("get pool error: %s", err.Error())
	//}
	//
	//c := conn.(net.Conn)

	//switch {
	//case data.Name() == "get":
	//	c.Write(task.reqeust.data)
	//case data.Name() == "set":
	//	c.Write(task.reqeust.data)
	//default:
	//	log.Println("no case ,", data.Name())
	//
	//}
	//write(ch, task, c)
}

func write(ch chan *task, task *task, c *redis.Conn) {
	defer c.Close()
	resp, _ := c.Decode()
	task.response = resp
	//redisBuffer := make([]byte, 100)
	//n, readerr := c.Read(redisBuffer)
	//
	//if readerr != nil {
	//	log.Errorln("read redis error: %s", readerr.Error())
	//}
	//
	//rr := bytes.NewReader(redisBuffer[:n])
	//cmd, _ := ReadCommand(rr)
	//resp := &Resp{
	//	data: redisBuffer[:n],
	//	cmd:  cmd,
	//}
	//task.response = resp
	//task.wg.Done()
	//ch <- task
}

func loopReadRedis(task *task) {

}

//TODO
func (rz *Router) CheckSlot(s *Slot) {

	switch {
	case s.status == MIGRATE:
		log.Println("MIGRATE")
	case s.status == OFFLINE:
		log.Println("OFFLINE")
	case s.status == ONLINE:
		log.Println("ONLINE")

	}

}
