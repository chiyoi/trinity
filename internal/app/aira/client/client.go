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
	trinityUrl string
	auth       string

	serviceUrl        string
	redisKeyListeners string
)

func init() {
	var err error
	if trinityUrl, err = config.Get[string]("TrinityURL"); err != nil {
		logs.Fatal(err)
	}

	if auth, err = config.Get[string]("TrinityAccessToken"); err != nil {
		logs.Fatal(err)
	}

	if serviceUrl, err = config.Get[string]("ServiceURL"); err != nil {
		logs.Fatal(err)
	}
	if redisKeyListeners, err = config.Get[string]("RedisKeyListeners"); err != nil {
		logs.Fatal(err)
	}
}

func PostMessage(a ...any) {
	if _, err := trinity.PostMessage(trinityUrl, auth, a...); err != nil {
		logs.Error("aira:", err)
		return
	}
}

func ErrorCallback(err error) {
	logs.Error("aira:", err)
	PostMessage(fmt.Sprintf("何処か間違ったような…[%s]", err))
}

func CacheFile(data []byte) (sasUrl string, err error) {
	return trinity.CacheFile(trinityUrl, auth, data)
}

func RegisterListener(rdb *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rdb.SAdd(ctx, redisKeyListeners, serviceUrl)
}

func SyncWorker(chanTimestamp chan int64, rdb *redis.Client) {
	timestamp := time.Now().Unix()
	for {
		select {
		case timestamp = <-chanTimestamp:
			continue
		case <-time.After(time.Second * 10):
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		b := rdb.SIsMember(ctx, redisKeyListeners, serviceUrl)
		if b.Err() != nil || b.Val() {
			continue
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		rdb.SAdd(ctx, redisKeyListeners, serviceUrl)
		cancel()

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