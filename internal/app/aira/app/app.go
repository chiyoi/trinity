package app

import (
	"time"

	"github.com/chiyoi/trinity/internal/pkg/logs"
	atmt2 "github.com/chiyoi/trinity/pkg/atmt"
)

func LogEvent(baseHandler atmt2.Handler, chanTimestamp chan<- int64) atmt2.Handler {
	return atmt2.HandlerFunc(func(ev atmt2.Event) {
		chanTimestamp <- time.Now().Unix()
		logs.Info("イベントが来たよ。")
		baseHandler.ServeEvent(ev)
		logs.Info("処理完了したよ。")
	})
}
