package trinity

import (
	"github.com/chiyoi/trinity/pkg/atmt/message"
)

type Request struct {
	Action Action `json:"action"`
	Data   any    `json:"data"`
}

type Action uint8

const (
	ActionPostMessage Action = iota
	ActionGetMessage
	ActionQueryMessageIdsTimeRange

	ActionCacheFile
)

func (act Action) String() string {
	switch act {
	case ActionPostMessage:
		return "post message"
	case ActionGetMessage:
		return "get message"
	case ActionQueryMessageIdsTimeRange:
		return "query message ids time range"
	case ActionCacheFile:
		return "cache file"
	default:
		return "invalid action"
	}
}

type ReqDataPostMessage struct {
	Message message.Message `json:"message"`
}

type ReqDataGetMessage struct {
	Id string `json:"id"`
}

type ReqDataQueryMessageTimeRange struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

type ReqDataCacheFile struct {
	Md5SumHex string `json:"md5_sum_hex"`
}

type Response[Data any] struct {
	StatusCode StatusCode `json:"status_code"`
	Data       Data       `json:"data"`
}

type StatusCode uint8

const (
	StatusOK StatusCode = iota
	StatusFailed
)

type RespDataPostMessage struct {
	MessageId string `json:"message_id"`
}

type RespDataGetMessage struct {
	Time      int64           `json:"time"`
	User      string          `json:"user"`
	MessageId string          `json:"message_id"`
	Message   message.Message `json:"message"`
}

type RespDataQueryMessageTimeRange struct {
	Ids []string `json:"ids"`
}

type RespDataCacheFile struct {
	SasURL string `json:"sas_url"`
}
