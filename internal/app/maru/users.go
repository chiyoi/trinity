package maru

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/onebot"
	"github.com/chiyoi/trinity/pkg/trinity"
	"github.com/go-redis/redis/v8"
)

var (
	redisKeyUsersLoggedIn string
)

func init() {
	var err error
	if redisKeyUsersLoggedIn, err = GetConfig[string]("RedisKeyUsersLoggedIn"); err != nil {
		logs.Fatal("maru:", err)
	}
}

func Login(rdb *redis.Client, id onebot.UserId, user, passwd string) {
	auth := trinity.CreateAuthorization(user, passwd)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rdb.HSet(ctx, redisKeyUsersLoggedIn, id, auth)
}

func GetAuthFromLoggedIn(rdb *redis.Client, id onebot.UserId) (auth string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return rdb.HGet(ctx, redisKeyUsersLoggedIn, strconv.Itoa(int(id))).Val()
}

func GetLoggedInList(rdb *redis.Client) (ids []onebot.UserId, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get logged in list: %w", err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	for _, id := range rdb.HKeys(ctx, redisKeyUsersLoggedIn).Val() {
		var iid int
		if iid, err = strconv.Atoi(id); err != nil {
			return
		}
		ids = append(ids, onebot.UserId(iid))
	}
	return
}
