package handlers

import (
	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func LogMessage(h atmt.Handler) atmt.Handler {
	return atmt.HandlerFunc(func(resp *atmt.Message, post atmt.Message) {
		logs.Infof("here comes a message~ [%s]", post.Type)
		h.ServeMessage(resp, post)
	})
}
