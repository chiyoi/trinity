package atmt

import (
	"github.com/chiyoi/trinity/pkg/atmt/message"
)

type Request struct {
	Time      int64           `json:"time"`
	User      string          `json:"user"`
	MessageId string          `json:"message_id"`
	Message   message.Message `json:"message"`
}

type Response struct {
	StatusCode StatusCode
}

type StatusCode uint8

const (
	StatusOK StatusCode = iota
	StatusFailed
)
