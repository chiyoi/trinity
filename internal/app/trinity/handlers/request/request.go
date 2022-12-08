package request

import (
	"encoding/json"

	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

func Request() (atmt.Matcher, atmt.Handler) {
	return atmt.Matcher{
		Match: func(msg atmt.Message) bool {
			return msg.Type == atmt.MessageRequest
		},
	}, handler()
}

func handler() atmt.Handler {
	return atmt.HandlerFunc(func(resp *atmt.Message, post atmt.Message) {
		var req atmt.DataRequest[trinity.Action]
		if err := json.Unmarshal(post.Data, &req); err != nil {
			logs.Warning("cannot parse request.")
			atmt.Error(resp, atmt.StatusBadRequest)
			return
		}
		switch req.Action {
		case trinity.ActionPostMessage:
			handlePostMessage(resp, req, post.Content)
		case trinity.ActionGetMessage:
		case trinity.ActionQueryMessageIdsTimeRange:
		case trinity.ActionCacheFile:
		case trinity.ActionVerifyAuthorization:
		default:
			logs.Warning("invalid action.")
			atmt.Error(resp, atmt.StatusBadRequest)
			return
		}
	})
}
