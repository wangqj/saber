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

	pool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:16379")
			if err != nil {
				log.Println("00000000", err)
				return nil, err
			}
			redis.DialConnectTimeout(3 * time.Second)
			return c, nil
		},
		MaxActive:       100,
		MaxIdle:         100,
		Wait:            true,
		MaxConnLifetime: 600 * time.Second,
	}
	//time.Sleep(10*time.Second)
	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Println("active count=", pool.ActiveCount())
		}
	}()

	t1 := time.Now()
	for i := 0; i < 1000000; i++ {
		go func(pool *redis.Pool, i int) {
			wg.Add(1)
			c := pool.Get()
			// pool.Close()
			defer c.Close()
			r, e := redis.String(c.Do("get", i))
			if e != nil {
				fmt.Println("1-------------", e, i)
			} else {
				fmt.Println("2------------", r, i)
				if i == 0 {
					fmt.Println("first cost time=", time.Since(t1))
				}
			}
			wg.Done()
		}(pool, i)
	}
	wg.Wait() // 等待
	t2 := time.Since(t1)
	fmt.Println("cost time=", t2)
	time.Sleep(1000 * time.Second)
}
