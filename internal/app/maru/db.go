package maru

import "github.com/go-redis/redis/v8"

func OpenRedis() (rdb *redis.Client, err error) {
	opt, err := GetConfig[*redis.Options]("RedisOptions")
	if err != nil {
		return
	}
	rdb = redis.NewClient(opt)
	return
}
