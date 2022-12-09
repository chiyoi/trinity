package handlers

import (
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func LogMessage(h atmt.Handler) atmt.Handler {
	return atmt.HandlerFunc(func(resp *atmt.Message, post atmt.Message) {
		logs.Info("here comes a message~", post.Type)
		h.ServeMessage(resp, post)
	})
}
