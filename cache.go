package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"strings"
	"sync"
	"time"
)

var cache = map[string]Cache{}
var cacheMutex = sync.Mutex{}

type Cache struct {
	body []byte
	ttl  time.Time
}

type CacheProxy struct {
	key string
	ttl time.Duration
}

func EvictCacheHandler(c *fiber.Ctx) error {
	return nil
}

func NewCacheProxy(key string, ttl time.Duration) CacheProxy {
	return CacheProxy{
		key: key,
		ttl: ttl,
	}
}

func (p CacheProxy) Accept(key string) bool {
	return p.key == key
}
func (p CacheProxy) Proxy(c *fiber.Ctx) error {
	path := c.Path()
	key := c.Params("key")

	if r, ok := cache[path]; ok && r.ttl.After(time.Now()) {
		c.Response().SetBody(r.body)
		c.Response().Header.Add("cache-control", fmt.Sprintf("max-age:%d", p.ttl/time.Second))
		return nil
	}

	url := "https://mocki.io/" + strings.TrimPrefix(path, "/"+key+"/")
	fmt.Printf("Http Request Redirecting to %s \n", url)

	if err := proxy.Do(c, url); err != nil {
		return err
	}

	ch := Cache{
		body: c.Response().Body(),
		ttl:  time.Now().Add(p.ttl),
	}

	cacheMutex.Lock()
	cache[path] = ch
	cacheMutex.Unlock()

	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
