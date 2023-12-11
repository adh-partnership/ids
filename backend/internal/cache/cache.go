package cache

import (
	"context"
	"errors"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	ErrorKeyNotFound = errors.New("key not found")
)

type Cache struct {
	ctx   context.Context
	r     *redis.Client
	inMem map[string]interface{}
	mut   sync.RWMutex
}

func NewCache(ctx context.Context, r *redis.Client) *Cache {
	return &Cache{
		ctx:   ctx,
		r:     r,
		inMem: make(map[string]interface{}),
	}
}

func (c *Cache) Get(key string) (interface{}, error) {
	if c.r == nil {
		c.mut.RLock()
		defer c.mut.RUnlock()
		if val, ok := c.inMem[key]; ok {
			return val, nil
		} else {
			return nil, ErrorKeyNotFound
		}
	}

	val, err := c.r.Get(c.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrorKeyNotFound
		}
		return nil, err
	}

	return val, nil
}

func (c *Cache) Set(key string, val interface{}) error {
	if c.r == nil {
		c.mut.Lock()
		defer c.mut.Unlock()
		c.inMem[key] = val
		return nil
	}

	err := c.r.Set(c.ctx, key, val, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Delete(key string) error {
	if c.r == nil {
		c.mut.Lock()
		defer c.mut.Unlock()
		delete(c.inMem, key)
		return nil
	}

	err := c.r.Del(c.ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
