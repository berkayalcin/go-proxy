package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sync"
	"time"
)

var counter = map[string]*Limit{}
var mutex = sync.Mutex{}

type Limit struct {
	count int
	ttl   time.Time
}

type LimitProxy struct {
	key   string
	limit int
	ttl   time.Duration
}

func NewLimitProxy(key string, limit int, ttl time.Duration) LimitProxy {
	return LimitProxy{
		key:   key,
		limit: limit,
		ttl:   ttl,
	}
}

func (p LimitProxy) Accept(key string) bool {
	return p.key == key
}

func (p LimitProxy) Proxy(c *fiber.Ctx) error {
	path := c.Path()

	if r, ok := counter[path]; ok && r.count >= p.limit && r.ttl.After(time.Now()) {
		c.Response().SetStatusCode(429)
		fmt.Printf("Request Limit Reached for %s\n", path)
		return nil
	} else if !ok {
		mutex.Lock()
		counter[path] = &Limit{
			count: 0,
			ttl:   time.Now().Add(p.ttl),
		}
		mutex.Unlock()
	}
	err := c.SendString("Go Türkiye - 103 Http Package")
	if err != nil {
		return err
	}

	mutex.Lock()
	counter[path].count++
	mutex.Unlock()
	return nil
}
