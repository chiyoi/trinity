package db

import (
	"sync"

	"github.com/chiyoi/trinity/internal/app/aira/config"
	"github.com/go-redis/redis/v8"
)

func OpenRedis() (rdb *redis.Client, err error) {
	opt, err := config.GetErr[*redis.Options]("RedisOptions")
	if err != nil {
		return
	}
	rdb = redis.NewClient(opt)
	return
}

var pool struct {
	rdb *redis.Client
	mu  sync.RWMutex
}

func SetDB(r *redis.Client) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if r != nil {
		pool.rdb = r
	}
}
func GetDB() (r *redis.Client) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	return pool.rdb
}
