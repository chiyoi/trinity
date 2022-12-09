package client

import (
	"context"
	"fmt"
	"time"

	"github.com/chiyoi/trinity/internal/app/aira/config"
	"github.com/chiyoi/trinity/internal/app/aira/db"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

var (
	serviceURL = config.Get[string]("ServiceURL")

	trinityURL = config.Get[string]("TrinityURL")
	auth       = config.Get[string]("TrinityAccessToken")

	redisKeyListeners = config.Get[string]("RedisKeyListeners")

	dbTimeout = time.Second * 10
)

func PostMessage(v ...any) {
	logPrefix := "post message:"
	if _, _, err := trinity.Request[trinity.ArgsPostMessage, trinity.ValuesPostMessage](
		trinityURL,
		trinity.ActionPostMessage,
		trinity.ArgsPostMessage{
			Auth: auth,
		},
		atmt.FormatContent(v...),
	); err != nil {
		logs.Error(logPrefix, err)
		return
	}
}

func CacheBlob(b []byte) (url string, err error) {
	return trinity.CacheBlob(trinityURL, auth, b)
}

var (
	errRdbNotSet = fmt.Errorf("rdb not set")
)

func CheckListener() (ok bool, err error) {
	rdb := db.GetDB()
	if rdb == nil {
		err = errRdbNotSet
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	cmd := rdb.SIsMember(ctx, redisKeyListeners, serviceURL)
	return cmd.Val(), cmd.Err()
}

func RegisterListener() (err error) {
	rdb := db.GetDB()
	if rdb == nil {
		return errRdbNotSet
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	return rdb.SAdd(ctx, redisKeyListeners, serviceURL).Err()
}

func RemoveListener() (err error) {
	rdb := db.GetDB()
	if rdb == nil {
		return errRdbNotSet
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	return rdb.SRem(ctx, redisKeyListeners, serviceURL).Err()
}
