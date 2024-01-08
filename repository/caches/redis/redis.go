package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

var (
	ErrDuplicateKey = errors.New("duplicate key")
	ErrNotFound     = errors.New("not found")
)

type Cache struct {
	rdb *redis.Client
}

func (c *Cache) Close() error {
	return c.rdb.Close()
}

func NewRedisCache(client *redis.Client) *Cache {
	result := Cache{rdb: client}
	return &result
}

func asInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
