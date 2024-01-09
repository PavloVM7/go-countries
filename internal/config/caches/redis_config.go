package caches

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Db       int
	Password string
}

func (r *RedisConfig) Read() {
	// set defaults
	r.Host = "127.0.0.1"
	r.Port = 6379
	r.Db = 0
	r.Password = ""
	// ToDo: load environment variables
}
func (r *RedisConfig) address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
func (r *RedisConfig) RedisOptions() redis.Options {
	result := redis.Options{
		Addr:     r.address(),
		DB:       r.Db,
		Password: r.Password,
	}
	return result
}
