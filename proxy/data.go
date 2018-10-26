package proxy

import (
	log "github.com/sirupsen/logrus"
	"time"
	"saber/proxy/redis"
)

type data struct {
	in chan *task
}

var ch = make(chan *task, 20480)
var dd = &data{in: ch}

func NewData() (*data) {
	return dd
}

func (d *data) AddData(t *task) {
	d.in <- t
}

func (d *data) Run(rz *Router) {
	for i := 0; i < 8; i++ {
		d.loopWrite(rz)
	}
}

func (d *data) GetConn(rz *Router) (*redis.Conn, chan *task) {
	c, err := redis.DialTimeout(":6381", time.Second*5,
		128*1024,
		128*1024)
	if err != nil {
		log.Errorln("HandleRequest 1resp error: %s", err.Error())
	}
	t := make(chan *task, 20480)
	go d.loopRead(t, c)

	return c, t
}

func (d *data) loopWrite(rz *Router) {
	c, t := d.GetConn(rz)
	p := c.FlushEncoder()
	p.MaxInterval = time.Millisecond
	p.MaxBuffered = cap(ch) / 2

	for r := range d.in {
		err := p.EncodeMultiBulk(r.Multi)
		if err != nil {
			log.Errorln("HandleRequest 2resp error: %s", err.Error())
		}
		if err := p.Flush(len(d.in) == 0); err != nil {
			log.Errorln("Flush 2resp error: %s", err.Error())
		}
		t <- r
	}
}

func (d *data) loopRead(rc chan *task, c *redis.Conn) {
	defer func() {
		c.Close()
	}()
	for r := range rc {
		resp, err := c.Decode()
		if err != nil {
			log.Errorln("read redis error: %s", err.Error())
		}
		r.response = resp
		r.wg.Done()
	}
}
