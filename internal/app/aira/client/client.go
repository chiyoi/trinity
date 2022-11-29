package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/chiyoi/trinity/internal/app/aira/config"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/trinity"
)

var (
	serviceURL = config.Get[string]("ServiceURL")

	trinityUrl = config.Get[string]("TrinityURL")
	auth       = config.Get[string]("TrinityAccessToken")

	redisKeyListeners = config.Get[string]("RedisKeyListeners")
)

func PostMessage(a ...any) {
	if _, err := trinity.PostMessage(trinityUrl, auth, a...); err != nil {
		logs.Error("aira:", err)
		return
	}
}

func CacheFile(data []byte) (sasUrl string, err error) {
	return trinity.CacheFile(trinityUrl, auth, data)
}

func RegisterListener(rdb *redis.Client) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return rdb.SAdd(ctx, redisKeyListeners, serviceURL).Err()
}

func EventSynchronizer(timestampChannel chan int64, rdb *redis.Client) (err error) {
	timestamp := time.Now().Unix()
	for {
		select {
		case timestamp = <-timestampChannel:
			continue
		case <-time.After(time.Second * 10):
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		cmd := rdb.SIsMember(ctx, redisKeyListeners, serviceURL)
		if err = cmd.Err(); err != nil {
			return
		}
		if cmd.Val() {
			continue
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err = rdb.SAdd(ctx, redisKeyListeners, serviceURL).Err(); err != nil {
			return
		}

		ids, err := trinity.QueryMessageIdsTimeRange(trinityUrl, auth, timestamp, time.Now().Unix())
		if err != nil {
			logs.Error("sync worker:", err)
			continue
		}
		for _, id := range ids {
			data, err := trinity.GetMessage(trinityUrl, auth, id)
			if err != nil {
				logs.Error("sync worker:", err)
				continue
			}

			req := atmt.Request{
				Time:      timestamp,
				User:      data.User,
				MessageId: id,
				Message:   data.Message,
			}
			b, err := json.Marshal(req)
			if err != nil {
				logs.Error("sync worker:", err)
				continue
			}

			resp, err := http.Post("localhost", "application/json", bytes.NewReader(b))
			if err != nil || resp.StatusCode != http.StatusOK {
				if err == nil {
					err = fmt.Errorf("api call error(%d %s)", resp.StatusCode, resp.Status)
				}
				logs.Error("sync worker:", err)
				continue
			}
		}
	}
}
