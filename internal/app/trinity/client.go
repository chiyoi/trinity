package trinity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/message"
)

var redisKeyListeners string

func init() {
	var err error
	if redisKeyListeners, err = GetConfig[string]("RedisKeyListeners"); err != nil {
		logs.Fatal("trinity:", err)
	}
}

func pushMsgToListeners(baseCtx context.Context, rdb *redis.Client, now int64, user string, messageId string, msg message.Message) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("push message to listeners: %w", err)
		}
	}()
	req := atmt.Request{
		Time:      now,
		User:      user,
		MessageId: messageId,
		Message:   msg,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(baseCtx, dbOperationTimeout)
	defer cancel()
	for _, l := range rdb.SMembers(ctx, redisKeyListeners).Val() {
		ctx, cancel := context.WithTimeout(baseCtx, time.Second*10)
		defer cancel()
		req, err1 := http.NewRequestWithContext(ctx, "POST", l, bytes.NewReader(data))
		if err1 != nil {
			logs.Error(err1)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err1 := http.DefaultClient.Do(req)
		if err1 != nil || resp.StatusCode != http.StatusOK {
			var logStr string
			if err != nil {
				logStr = err1.Error()
			} else {
				logStr = resp.Status
			}
			logs.Warning("push error:", logStr)
			ctx, cancel = context.WithTimeout(baseCtx, time.Second*10)
			defer cancel()
			rdb.SRem(ctx, redisKeyListeners, l)
			continue
		}
	}
	return
}
