package main

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
	"time"
	"sync"
	"log"
)

func main() {
	var wg sync.WaitGroup
	t1 := time.Now()
	pool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:16379")
			if err != nil {
				log.Println("00000000", err)
				return nil, err
			}
			return c, nil
		},
		MaxActive:       100,
		MaxIdle:         10,
		Wait:            true,
		MaxConnLifetime: 60 * time.Second,
	}
	//time.Sleep(10*time.Second)
	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Println("active count=", pool.ActiveCount())
		}
	}()
	for i := 0; i < 10000; i++ {
		//go func(pool *redis.Pool,i int) {
		wg.Add(1)
		c := pool.Get()
		// pool.Close()
		fmt.Println("0--------", c, i)
		r, e := c.Do("get", i)
		r2, e2 := c.Do("get", i)
		c.Do("get", i)
		c.Do("get", i)
		c.Do("get", i)
		c.Do("get", i)
		c.Do("get", i)
		c.Do("get", i)
		c.Do("get", i)
		c.Do("get", i)
		c.Do("get", i)
		c.Close()
		if e != nil {
			fmt.Println("1-------------", e, i)
			fmt.Println("1-------------", e2, i)
		} else {
			fmt.Println("2------------", r, i)
			fmt.Println("2------------", r2, i)

		}
		fmt.Println("3--------", c, i)
		//}(pool,i)
	}
	wg.Wait() // 等待
	t2 := time.Since(t1)
	fmt.Println("cost time=", t2)
	time.Sleep(100 * time.Second)
}
