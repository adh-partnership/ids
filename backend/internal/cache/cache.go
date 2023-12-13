package cache

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/adh-partnership/ids/backend/internal/redis"
)

const (
	NoExpiration      = -1
	DefaultExpiration = 0
)

var (
	ErrorKeyNotFound = errors.New("key not found")
)

type Item struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	defaultExpiration time.Duration
	ctx               context.Context
	redis             *redis.Redis
	inMem             map[string]Item
	mut               sync.RWMutex
	janitor           *janitor
}

func NewCache(ctx context.Context, r *redis.Redis, defaultExpiration time.Duration, janitorInterval time.Duration) *Cache {
	if janitorInterval == 0 {
		janitorInterval = 2 * time.Minute
	}

	c := &Cache{
		ctx:   ctx,
		redis: r,
		inMem: make(map[string]Item),
	}

	// We only need the janitor if we are not using redis. Redis will handle item expiration on its own.
	if r == nil {
		j := newJanitor(janitorInterval)
		j.Run(c)
		runtime.SetFinalizer(c, stopJanitor)
	}

	return c
}

func (c *Cache) Get(key string) (interface{}, error) {
	if c.redis == nil {
		c.mut.RLock()
		defer c.mut.RUnlock()
		if val, ok := c.inMem[key]; ok {
			return val.Value, nil
		} else {
			return nil, ErrorKeyNotFound
		}
	}

	val, err := c.redis.Get(c.ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, ErrorKeyNotFound
		}
		return nil, err
	}

	return val, nil
}

func (c *Cache) SetDefault(key string, val interface{}) error {
	return c.Set(key, val, c.defaultExpiration)
}

func (c *Cache) Set(key string, val interface{}, d time.Duration) error {
	if c.redis == nil {
		var exp int64

		if d == DefaultExpiration {
			d = c.defaultExpiration
		}

		if d > 0 {
			exp = time.Now().Add(d).UnixNano()
		}

		c.mut.Lock()
		c.inMem[key] = Item{
			Value:      val,
			Expiration: exp,
		}
		c.mut.Unlock()

		return nil
	}

	var exp time.Duration
	if d == DefaultExpiration {
		exp = c.defaultExpiration
	} else if d == NoExpiration {
		exp = 0
	}

	err := c.redis.Set(c.ctx, key, val, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Delete(key string) error {
	if c.redis == nil {
		c.mut.Lock()
		defer c.mut.Unlock()
		delete(c.inMem, key)
		return nil
	}

	err := c.redis.Del(c.ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mut.Lock()
	for k, v := range c.inMem {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.inMem, k)
		}
	}
	c.mut.Unlock()
}
