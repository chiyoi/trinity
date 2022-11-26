package aira

import (
	"github.com/chiyoi/trinity/internal/app/aira/config"
	"github.com/go-redis/redis/v8"
)

func OpenRedis() (rdb *redis.Client, err error) {
	opt, err := config.Get[*redis.Options]("RedisOptions")
	if err != nil {
		return
	}
	rdb = redis.NewClient(opt)
	return
}
