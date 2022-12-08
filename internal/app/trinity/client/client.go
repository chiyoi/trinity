package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/chiyoi/trinity/internal/app/trinity/db"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

var (
	dbTimeout   = time.Second * 10
	pushTimeout = time.Second * 10
)

func PushMessageToListeners(msg atmt.Message) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("push event to listeners: %w", err)
		}
	}()
	bg := context.Background()
	ls := config.Get[string]("RedisKeyListeners")
	rdb, _ := db.GetDB()
	if rdb == nil {
		err = errors.New("rdb not set")
		return
	}
	ctx, cancel := context.WithTimeout(bg, dbTimeout)
	defer cancel()
	cmd := rdb.SMembers(ctx, ls)
	if err = cmd.Err(); err != nil {
		return
	}
	for _, l := range cmd.Val() {
		ctx, cancel := context.WithTimeout(bg, pushTimeout)
		defer cancel()
		if err = atmt.PostCtx(ctx, l, msg); err != nil {
			if _, ok := err.(*url.Error); !ok {
				return
			}
			logs.Warning(fmt.Sprintf("cannot push to listener %s: %s", l, err))
			ctx, cancel = context.WithTimeout(bg, dbTimeout)
			defer cancel()
			if err = rdb.SRem(ctx, ls, l).Err(); err != nil {
				return
			}
		}
	}
	return
}
