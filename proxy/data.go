package proxy

import (
	log "github.com/sirupsen/logrus"
	"time"
	"saber/proxy/redis"
	"sync"
	"saber/utils"
)

type data struct {
	input chan *task
	inner chan *task
}

var once sync.Once
var instance *data

func NewData(o *utils.Option) *data {
	log.Println("init data ", o.SessionBuffer)
	once.Do(func() {
		instance = &data{
			input: make(chan *task, o.SessionBuffer),
			inner: make(chan *task, o.DataBuffer),
		}
	})
	return instance
}

func GetData() (*data) {
	return instance
}

func (d *data) Run(rz *Router) {

	for i := 0; i < 8; i++ {
		d.loopWrite(rz)
	}
}

func (d *data) GetConn(rz *Router) (*redis.Conn) {
	c, err := redis.DialTimeout(":6381", time.Second*5,
		128*1024,
		128*1024)
	if err != nil {
		log.Errorln("GetConn error: %s", err.Error())
	}
	go d.loopRead(c)

	return c
}

func (d *data) loopWrite(rz *Router) {
	c := d.GetConn(rz)
	p := c.FlushEncoder()
	p.MaxInterval = time.Millisecond
	p.MaxBuffered = cap(d.inner) / 2

	for r := range d.input {
		err := p.EncodeMultiBulk(r.Multi)
		if err != nil {
			log.Errorln("EncodeMultiBulk error: %s", err.Error())
		}
		if err := p.Flush(len(d.input) == 0); err != nil {
			log.Errorln("Flush resp error: %s", err.Error())
		}
		d.inner <- r
	}
}

func (d *data) loopRead(c *redis.Conn) {
	defer func() {
		c.Close()
	}()
	for r := range d.inner {
		resp, err := c.Decode()
		if err != nil {
			log.Errorln("read redis error: %s", err.Error())
		}
		r.response = resp
		r.wg.Done()
	}
}
