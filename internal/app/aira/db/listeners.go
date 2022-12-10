package db

import (
	"context"

	"github.com/chiyoi/trinity/internal/app/aira/config"
	"github.com/go-redis/redis/v8"
)

func RdbKeyListeners() (rdb *redis.Client, key string, err error) {
	rdb = GetDB()
	if rdb == nil {
		err = errRdbNotSet
		return
	}
	key = config.Get[string]("KeyListeners")
	return
}

func CheckListener(serviceURL string) (ok bool, err error) {
	rdb, key, err := RdbKeyListeners()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	cmd := rdb.SIsMember(ctx, key, serviceURL)
	return cmd.Val(), cmd.Err()
}

func RegisterListener(serviceURL string) (err error) {
	rdb, key, err := RdbKeyListeners()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	return rdb.SAdd(ctx, key, serviceURL).Err()
}

func RemoveListener(serviceURL string) (err error) {
	rdb, key, err := RdbKeyListeners()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	return rdb.SRem(ctx, key, serviceURL).Err()
}
