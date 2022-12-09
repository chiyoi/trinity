package trinity

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Action uint8

const (
	ActionPostMessage Action = iota + 1
	ActionGetMessage
	ActionQueryMessageIdsLatestCount

	ActionGetBlobCacheURL
	ActionVerifyAuthorization
)

func (act Action) String() (str string) {
	defer func() {
		str = fmt.Sprintf("Action(%s)", str)
	}()
	switch act {
	case ActionPostMessage:
		return "post message"
	case ActionGetMessage:
		return "get message"
	case ActionQueryMessageIdsLatestCount:
		return "query message-ids time range"
	case ActionGetBlobCacheURL:
		return "get blob cache url"
	case ActionVerifyAuthorization:
		return "verify authorization"
	default:
		return "invalid action"
	}
}

type void = struct{}

type MessageID = primitive.ObjectID

type ArgsPostMessage struct {
	Auth string `json:"auth"`
}
type ValuesPostMessage struct {
	MessageID MessageID `json:"message_id"`
}

type ArgsGetMessage struct {
	Auth string    `json:"auth"`
	ID   MessageID `json:"id"`
}
type ValuesGetMessage struct {
	Sender    string    `json:"sender"`
	MessageID MessageID `json:"message_id"`
}

type ValuesQueryMessageIds struct {
	Ids []MessageID `json:"ids"`
}
type ArgsQueryMessageIdsLatestCount struct {
	Auth  string `json:"auth"`
	Count int    `json:"count"`
}
type ValuesQueryMessageIdsLatestCount = ValuesQueryMessageIds

type ArgsGetBlobCacheURL struct {
	Auth     string `json:"auth"`
	BlobName string `json:"blob_name"`
}
type ValuesGetBlobCacheURL struct {
	SasURL string `json:"sas_url"`
}

type ArgsVerifyAuthorization struct {
	Auth string `json:"auth"`
}
type ValuesVerifyAuthorization struct {
	Pass bool `json:"pass"`
}
