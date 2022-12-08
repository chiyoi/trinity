package trinity

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Action uint8

const (
	ActionPostMessage Action = iota + 1
	ActionGetMessage
	ActionQueryMessageIdsTimeRange

	ActionCacheFile
	ActionVerifyAuthorization
)

func (act Action) String() (str string) {
	defer func() {
		str = fmt.Sprintf("Action(%s)", str)
	}()
	switch act {
	case ActionPostMessage:
		return "post message"
	default:
		return "invalid action"
	}
}

type MessageId = primitive.ObjectID
type ArgsPostMessage struct {
	Sender string `json:"sender"`
	Auth   string `json:"auth"`
}

type ArgsGetMessage struct {
	Auth string `json:"auth"`
	Id   string `json:"id"`
}

type ArgsQueryMessageTimeRange struct {
	Auth string `json:"auth"`
	From int64  `json:"from"`
	To   int64  `json:"to"`
}

type ArgsCacheFile struct {
	Auth         string `json:"auth"`
	Sha256SumHex string `json:"sha256_sum_hex"`
}

type ArgsVoid struct{}

type RequestBuilder[Args any] struct {
	Action Action `json:"action"`
	Args   Args   `json:"args"`
}
