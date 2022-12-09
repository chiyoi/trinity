package handlers

import (
	"fmt"

	"github.com/chiyoi/trinity/internal/app/aira/client"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func LogMessage(h atmt.Handler) atmt.Handler {
	return atmt.HandlerFunc(func(resp *atmt.Message, post atmt.Message) {
		logs.Info("メッセージが来たよ〜", post.Type)
		h.ServeMessage(resp, post)
	})
}

func Error(err error) {
	client.PostMessage(fmt.Sprintf("何処か間違ったような…(%s)", err))
}
