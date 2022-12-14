package db

import (
	"context"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/go-redis/redis/v8"
)

func RdbKeyListeners() (rdb *redis.Client, key string, err error) {
	rdb, _ = GetDB()
	if rdb == nil {
		err = errRdbNotSet
		return
	}
	key = config.Get[string]("KeyListeners")
	return
}

func GetListeners() (ls []string, err error) {
	rdb, key, err := RdbKeyListeners()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	cmd := rdb.SMembers(ctx, key)
	return cmd.Val(), cmd.Err()
}

func RemoveListener(l string) (err error) {
	rdb, key, err := RdbKeyListeners()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	return rdb.SRem(ctx, key, l).Err()
}
