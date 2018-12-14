package proxy

import (
	log "github.com/sirupsen/logrus"
	"time"
	"saber/proxy/redis"
	"sync"
	"saber/utils"
)

type Processor struct {
	mu    sync.Mutex
	conn  *redis.Conn
	input chan *task
	inner chan *task
}

func NewProcessor(conn *redis.Conn) *Processor {
	o := utils.Cfg
	log.Println("init Processor ", o.SessionBuffer)
	instance := &Processor{
		input: make(chan *task, o.SessionBuffer),
		inner: make(chan *task, o.DataBuffer),
		conn:  conn,
	}
	return instance
}

func (d *Processor) Start() {
	d.mu.Lock()
	defer d.mu.Unlock()
	go d.loopRead()
	d.loopWrite()
}

//TODO
func (d *Processor) Stop() {
	//close(d.input)
}

func (d *Processor) loopWrite() {
	p := d.conn.FlushEncoder()
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

func (d *Processor) loopRead() {
	defer func() {
		d.conn.Close()
	}()
	for r := range d.inner {
		resp, err := d.conn.Decode()
		if err != nil {
			log.Errorln("read redis error: %s", err.Error())
		}
		r.response = resp
		r.wg.Done()
	}
}
