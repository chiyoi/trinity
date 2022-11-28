package handlers

import (
	"fmt"
	"time"

	"github.com/chiyoi/trinity/internal/app/aira/client"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func LogEvent(baseHandler atmt.Handler, chanTimestamp chan<- int64) atmt.Handler {
	return atmt.HandlerFunc(func(ev atmt.Event) {
		chanTimestamp <- time.Now().Unix()
		logs.Info("イベントが来たよ。")
		baseHandler.ServeEvent(ev)
		logs.Info("処理完了したよ。")
	})
}

func ErrorCallback(err error) {
	logs.Error("aira:", err)
	client.PostMessage(fmt.Sprintf("何処か間違ったような…[%s]", err))
}
