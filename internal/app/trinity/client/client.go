package client

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/chiyoi/trinity/internal/app/trinity/db"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

const (
	pushTimeout = time.Second * 10
)

func PushMessageToListeners(msg atmt.Message) (err error) {
	logPrefix := "push message to listeners:"
	defer func() {
		if err != nil {
			err = fmt.Errorf("push message to listeners: %w", err)
		}
	}()
	ls, err := db.GetListeners()
	if err != nil {
		return
	}
	for _, l := range ls {
		go func(l string) {
			ctx, cancel := context.WithTimeout(context.Background(), pushTimeout)
			defer cancel()
			if err = atmt.PushCtx(ctx, l, msg); err != nil {
				if _, ok := err.(*url.Error); !ok {
					return
				}
				logs.Warning(logPrefix, fmt.Sprintf("push to listener %s: %s", l, err))
				if err = db.RemoveListener(l); err != nil {
					return
				}
			}
		}(l)
	}
	return
}
