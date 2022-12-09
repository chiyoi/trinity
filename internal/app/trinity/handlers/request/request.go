package request

import (
	"encoding/json"

	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

var Matcher = atmt.Matcher{
	Match: func(msg atmt.Message) bool {
		return msg.Type == atmt.MessageRequest
	},
}

type reqHandler = func(resp *atmt.Message, req atmt.DataRequest[trinity.Action])

func Handler() atmt.HandlerFunc {
	handleCacheFile, err := getBlobCacheURLHandler()
	if err != nil {
		logs.Fatal(err)
	}
	return func(resp *atmt.Message, post atmt.Message) {
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
			handleGetMessage(resp, req)
		case trinity.ActionQueryMessageIdsLatestCount:
			handleQueryMessageIdsLatestCount(resp, req)
		case trinity.ActionGetBlobCacheURL:
			handleCacheFile(resp, req)
		case trinity.ActionVerifyAuthorization:
			handleVerifyAuthorization(resp, req)
		default:
			logs.Warning("invalid action.")
			atmt.Error(resp, atmt.StatusBadRequest)
			return
		}
	}
}
