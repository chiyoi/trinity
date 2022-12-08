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

func PostMessage(a ...any) {
	logPrefix := "post message:"
	var content []atmt.Paragraph
	for _, aa := range a {
		switch t := aa.(type) {
		case atmt.Paragraph:
			content = append(content, t)
		default:
			content = append(content, atmt.Paragraph{
				Type: atmt.ParagraphText,
				Text: fmt.Sprint(t),
			})
		}
	}
	if err := trinity.PostMessage(trinityURL, auth, content); err != nil {
		logs.Error(logPrefix, err)
		return
	}
}

// func CacheFile(data []byte) (sasUrl string, err error) {
// 	return trinity.CacheFile(trinityUrl, auth, data)
// }

var (
	errRdbNotSet = fmt.Errorf("rdb not set")
)

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

// func EventSynchronizer(timestampChannel chan int64) (err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("event synchronizer: %w", err)
// 		}
// 	}()
// 	rdb, _ := db.GetDB()
// 	if rdb == nil {
// 		return errRdbNotSet
// 	}
// 	timestamp := time.Now().Unix()
// 	for {
// 		select {
// 		case timestamp = <-timestampChannel:
// 			continue
// 		case <-time.After(time.Second * 10):
// 		}

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 		defer cancel()
// 		cmd := rdb.SIsMember(ctx, redisKeyListeners, serviceURL)
// 		if err = cmd.Err(); err != nil {
// 			return
// 		}
// 		if cmd.Val() {
// 			continue
// 		}

// 		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
// 		defer cancel()
// 		if err = rdb.SAdd(ctx, redisKeyListeners, serviceURL).Err(); err != nil {
// 			return
// 		}

// 		ids, err := trinity.QueryMessageIdsTimeRange(trinityUrl, auth, timestamp, time.Now().Unix())
// 		if err != nil {
// 			logs.Error("sync worker:", err)
// 			continue
// 		}
// 		for _, id := range ids {
// 			data, err := trinity.GetMessage(trinityUrl, auth, id)
// 			if err != nil {
// 				logs.Error("sync worker:", err)
// 				continue
// 			}

// 			req := atmt.Request{
// 				Time:      timestamp,
// 				User:      data.User,
// 				MessageId: id,
// 				Message:   data.Message,
// 			}
// 			b, err := json.Marshal(req)
// 			if err != nil {
// 				logs.Error("sync worker:", err)
// 				continue
// 			}

// 			resp, err := http.Post("localhost", "application/json", bytes.NewReader(b))
// 			if err != nil || resp.StatusCode != http.StatusOK {
// 				if err == nil {
// 					err = fmt.Errorf("api call error(%d %s)", resp.StatusCode, resp.Status)
// 				}
// 				logs.Error("sync worker:", err)
// 				continue
// 			}
// 		}
// 	}
// }
