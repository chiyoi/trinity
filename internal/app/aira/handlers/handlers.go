package handlers

import (
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func LogEvent(baseHandler atmt.Handler) atmt.Handler {
	return atmt.HandlerFunc(func(resp *atmt.Message, post atmt.Message) {
		logs.Info("イベントが来たよ。")
		baseHandler.ServeMessage(resp, post)
		logs.Info("処理完了したよ。")
	})
}

// func ErrorCallback(err error) {
// 	logs.Error(err)
// 	client.PostMessage(fmt.Sprintf("何処か間違ったような…[%s]", err))
// }
