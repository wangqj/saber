package proxy

import (
	"sync"
	"saber/proxy/redis"
)

type Rst struct {
	request  *redis.Resp
	response *redis.Resp

	wg *sync.WaitGroup
}

type RstChan struct {
	lock sync.Mutex
	cond *sync.Cond

	data []*Rst
	buff []*Rst

	waits  int
	closed bool
}

func NewRequestChan() *RstChan {
	return NewRequestChanBuffer(0)
}

func NewRequestChanBuffer(n int) *RstChan {
	if n <= 0 {
		n = 1024
	}
	var ch = &RstChan{
		buff: make([]*Rst, n),
	}
	ch.cond = sync.NewCond(&ch.lock)
	return ch
}

func (c *RstChan) Close() {
	c.lock.Lock()
	if !c.closed {
		c.closed = true
		c.cond.Broadcast()
	}
	c.lock.Unlock()
}

func (c *RstChan) Buffered() int {
	c.lock.Lock()
	n := len(c.data)
	c.lock.Unlock()
	return n
}

func (c *RstChan) PushBack(r *Rst) int {
	c.lock.Lock()
	n := c.lockedPushBack(r)
	c.lock.Unlock()
	return n
}

func (c *RstChan) PopFront() (*Rst, bool) {
	c.lock.Lock()
	r, ok := c.lockedPopFront()
	c.lock.Unlock()
	return r, ok
}

func (c *RstChan) lockedPushBack(r *Rst) int {
	if c.closed {
		panic("send on closed chan")
	}
	if c.waits != 0 {
		c.cond.Signal()
	}
	c.data = append(c.data, r)
	return len(c.data)
}

func (c *RstChan) lockedPopFront() (*Rst, bool) {
	for len(c.data) == 0 {
		if c.closed {
			return nil, false
		}
		c.data = c.buff[:0]
		c.waits++
		c.cond.Wait()
		c.waits--
	}
	var r = c.data[0]
	c.data, c.data[0] = c.data[1:], nil
	return r, true
}

func (c *RstChan) IsEmpty() bool {
	return c.Buffered() == 0
}

func (c *RstChan) PopFrontAll(onRequest func(r *Rst) error) error {
	for {
		r, ok := c.PopFront()
		if ok {
			if err := onRequest(r); err != nil {
				return err
			}
		} else {
			return nil
		}
	}
}

func (c *RstChan) PopFrontAllVoid(onRequest func(r *Rst)) {
	c.PopFrontAll(func(r *Rst) error {
		onRequest(r)
		return nil
	})
}
