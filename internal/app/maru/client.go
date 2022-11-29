package maru

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/chiyoi/trinity/internal/app/maru/config"
	"github.com/chiyoi/trinity/pkg/onebot"
	"github.com/chiyoi/trinity/pkg/trinity"
	"github.com/chiyoi/websocket"
	"github.com/go-redis/redis/v8"
)

var (
	trinityUrl = config.Get[string]("TrinityURL")
	onebotUrl  = config.Get[string]("OnebotURL")

	redisKeyUsersLoggedIn = config.Get[string]("RedisKeyUsersLoggedIn")
	redisKeyListeners     = config.Get[string]("RedisKeyListeners")
)

func DialOnebotEventServer() (ws websocket.WebSocket, err error) {
	onebotEventUrl, err := config.GetErr[string]("OnebotEventURL")
	if err != nil {
		return
	}

	if ws, err = websocket.Dial(onebotEventUrl); err != nil {
		return
	}
	return
}

func Login(rdb *redis.Client, id onebot.UserId, user, passwd string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("login: %w", err)
		}
	}()
	auth := trinity.CreateAuthorization(user, passwd)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return rdb.HSet(ctx, redisKeyUsersLoggedIn, id, auth).Err()
}

func GetAuthFromLoggedIn(rdb *redis.Client, id onebot.UserId) (auth string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get auth from logged-in: %w", err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	cmd := rdb.HGet(ctx, redisKeyUsersLoggedIn, strconv.Itoa(int(id)))
	return cmd.Val(), cmd.Err()
}

func GetLoggedInList(rdb *redis.Client) (ids []onebot.UserId, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get logged-in list: %w", err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	cmd := rdb.HKeys(ctx, redisKeyUsersLoggedIn)
	if err = cmd.Err(); err != nil {
		return
	}
	for _, id := range cmd.Val() {
		var iid int
		if iid, err = strconv.Atoi(id); err != nil {
			return
		}
		ids = append(ids, onebot.UserId(iid))
	}
	return
}

func RegisterListener(rdb *redis.Client) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("register listener: %w", err)
		}
	}()
	serviceUrl, err := config.GetErr[string]("ServiceURL")
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return rdb.SAdd(ctx, redisKeyListeners, serviceUrl).Err()
}

func OnebotSendMsg(id onebot.UserId, a ...any) (err error) {
	_, err = onebot.SendMsg(onebotUrl, id, a...)
	return
}
