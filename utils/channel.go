package utils

import (
	"errors"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
	"fmt"
)

// PoolConfig 连接池相关配置
type PoolConfig struct {
	Name string
	//连接池中拥有的最小连接数
	InitialCap int
	//连接池中拥有的最大的连接数
	MaxCap int
	//生成连接的方法
	Factory func() (interface{}, error)
	//关闭连接的方法
	Close func(interface{}) error
	//检查连接是否有效的方法
	Ping func(interface{}) error
	//连接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
}

//channelPool 存放连接信息
type channelPool struct {
	name        string
	mu          sync.Mutex
	conns       chan *idleConn
	minCount    int
	maxCount    int
	activeCount int
	factory     func() (interface{}, error)
	close       func(interface{}) error
	ping        func(interface{}) error
	idleTimeout time.Duration
}

type idleConn struct {
	conn interface{}
	t    time.Time
}

//NewChannelPool 初始化连接
func NewChannelPool(poolConfig *PoolConfig) (Pool, error) {
	if poolConfig.InitialCap < 0 || poolConfig.MaxCap <= 0 || poolConfig.InitialCap > poolConfig.MaxCap {
		return nil, errors.New("invalid capacity settings")
	}
	if poolConfig.Factory == nil {
		return nil, errors.New("invalid factory func settings")
	}
	if poolConfig.Close == nil {
		return nil, errors.New("invalid close func settings")
	}

	c := &channelPool{
		name:        poolConfig.Name,
		conns:       make(chan *idleConn, poolConfig.MaxCap),
		minCount:    poolConfig.InitialCap,
		maxCount:    poolConfig.MaxCap,
		factory:     poolConfig.Factory,
		close:       poolConfig.Close,
		idleTimeout: poolConfig.IdleTimeout,
	}

	if poolConfig.Ping != nil {
		c.ping = poolConfig.Ping
	}

	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory()
		log.Println("pool stat:", c.name, ", current num is ", len(c.getConns()))
		if err != nil {
			c.Release()
			log.Errorln("factory is not able to fill the pool ", err)
			return nil, err
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}
	go report(c)
	go check(c)
	return c, nil
}

//getConns 获取所有连接
func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

//Get 从pool中取一个连接
func (c *channelPool) Get() (interface{}, error) {

	conns := c.getConns()
	if conns == nil {
		fmt.Println("if error ")
		return nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			//fmt.Println("get conn")
			if wrapConn == nil {
				fmt.Println("wrapConn error ")
				continue
			}
			//判断是否超时，超时则丢弃
			if timeout := c.idleTimeout; timeout > 0 && c.activeCount >= c.minCount {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该连接
					c.Close(wrapConn.conn)
					continue
				}
			}
			//判断是否失效，失效则丢弃，如果用户没有设定 ping 方法，就不检查
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					fmt.Println("conn is not able to be connected: ", err)
					continue
				}
			}
			c.mu.Lock()
			c.activeCount++
			c.mu.Unlock()
			return wrapConn.conn, nil
		default:
			if c.activeCount >= c.maxCount {
				//fmt.Println("max=",c.activeCount)
				//time.Sleep(100*time.Millisecond)
				continue
			}
			log.Println("default pool stat:", c.name, ", current num is ", len(c.getConns()))
			conn, err := c.factory()
			if err != nil {
				fmt.Println("default error ", err)
				return nil, err
			}
			c.mu.Lock()
			c.activeCount++
			c.mu.Unlock()
			return conn, nil
		}
	}
}

//Put 将连接放回pool中
func (c *channelPool) Put(conn interface{}) error {
	//fmt.Println("return ")
	if conn == nil {
		c.mu.Lock()
		c.activeCount--
		c.mu.Unlock()
		return errors.New("connection is nil. rejecting")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.activeCount++
	if c.conns == nil {
		return c.Close(conn)
	}

	select {
	case c.conns <- &idleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该连接
		return c.Close(conn)
	}
}

//Close 关闭单条连接
func (c *channelPool) Close(conn interface{}) error {
	c.mu.Lock()
	c.activeCount--
	c.mu.Unlock()
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.close(conn)
}

//Ping 检查单条连接是否有效
func (c *channelPool) Ping(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.ping(conn)
}

//Release 释放连接池中所有连接
func (c *channelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	closeFun := c.close
	c.close = nil
	c.activeCount = 0
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

//Len 连接池中已有的连接
func (c *channelPool) Len() int {
	return len(c.getConns())
}

func report(c *channelPool) {
	for {
		time.Sleep(10 * time.Second)
		log.Println("pool stat:", c.name, ", current num is ", len(c.getConns()))
	}

}

func check(c *channelPool) {
	for {
		time.Sleep(60 * time.Second)

		conns := c.getConns()
		if conns == nil {
			continue
		}

		select {
		case wrapConn := <-conns:
			f := false
			if wrapConn == nil {
				continue
			}
			//判断是否超时，超时则丢弃
			if timeout := c.idleTimeout; timeout > 0 && len(conns) >= c.minCount {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该连接
					c.Close(wrapConn.conn)
					f = true
					continue
				}
			}
			//判断是否失效，失效则丢弃，如果用户没有设定 ping 方法，就不检查
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					fmt.Println("conn is not able to be connected: ", err)
					f = true
					continue
				}
			}
			if !f {
				c.conns <- &idleConn{conn: wrapConn.conn, t: wrapConn.t}
			}
		}
	}
}
