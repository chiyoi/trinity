package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

var (
	simpleRequestTimeout = time.Second * 10
)

func init() {
	var err error
	if err != nil {
		logs.Fatal("trinity:", err)
	}
}

func PushEventToListeners(baseCtx context.Context, rdb *redis.Client, ev atmt.Event) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("push event to listeners: %w", err)
		}
	}()
	ls := config.Get[string]("RedisKeyListeners")

	ctx, cancel := context.WithTimeout(baseCtx, simpleRequestTimeout)
	defer cancel()
	cmd := rdb.SMembers(ctx, ls)
	if err = cmd.Err(); err != nil {
		return
	}
	for _, l := range cmd.Val() {
		ctx, cancel := context.WithTimeout(baseCtx, time.Second*10)
		defer cancel()
		if apiErr := atmt.CallApiCtx(ctx, l, atmt.Request{
			Time:      ev.Time.Unix(),
			User:      ev.User,
			MessageId: ev.MessageId,
			Message:   ev.Message,
		}); apiErr != nil {
			logs.Warning("push event error:", apiErr)

			ctx, cancel = context.WithTimeout(baseCtx, simpleRequestTimeout)
			defer cancel()
			if err = rdb.SRem(ctx, ls, l).Err(); err != nil {
				return
			}
			continue
		}
	}
	return
}
