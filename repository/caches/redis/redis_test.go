package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"os"
	"pm.com/go-countries/internal/config/caches"
	"testing"
	"time"
)

type redisTestSuit struct {
	suite.Suite
	rdb *redis.Client
}

func (s *redisTestSuit) SetupSuite() {
	var config caches.RedisConfig
	config.Read()
	redisOptions := config.RedisOptions()
	fmt.Printf("try connect to redis, config:%+v", config)
	s.rdb = redis.NewClient(&redisOptions)
	if s.rdb == nil {
		panic("couldn't connect to redis")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := s.rdb.Ping(ctx).Result()
	fmt.Println("redis ping result:", res)
	if err != nil {
		panic(err)
	}
}
func (s *redisTestSuit) TearDownSuite() {
	err := s.rdb.Close()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "close redis error:", err)
	} else {
		fmt.Println("redis connection closed")
	}
}

func (s *redisTestSuit) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	status, err := s.rdb.FlushAll(ctx).Result()
	fmt.Println("clear all databases", "result:", status, "err:", err)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(redisTestSuit))
}
