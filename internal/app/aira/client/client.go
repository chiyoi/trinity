package client

import (
	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/trinity/internal/app/aira/config"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

var (
	trinityURL = config.Get[string]("TrinityURL")
	auth       = config.Get[string]("TrinityAccessToken")
)

func PostMessage(v ...any) {
	if _, _, err := trinity.Request[trinity.ArgsPostMessage, trinity.ValuesPostMessage](
		trinityURL,
		trinity.ActionPostMessage,
		trinity.ArgsPostMessage{
			Auth: auth,
		},
		atmt.FormatContent(v...),
	); err != nil {
		logs.Error(err)
		return
	}
}

func CacheBlob(b []byte) (url string, err error) {
	return trinity.CacheBlob(trinityURL, auth, b)
}
