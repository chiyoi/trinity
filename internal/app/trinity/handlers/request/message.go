package request

import (
	"encoding/json"

	"github.com/chiyoi/trinity/internal/app/trinity/client"
	"github.com/chiyoi/trinity/internal/app/trinity/db"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

func handlePostMessage(resp *atmt.Message, req atmt.DataRequest[trinity.Action], content []atmt.Paragraph) {
	logPrefix := "handle post message:"
	var args trinity.ArgsPostMessage
	if err := json.Unmarshal(req.Args, &args); err != nil {
		logs.Warning(logPrefix, err)
		atmt.Error(resp, atmt.StatusBadRequest)
		return
	}
	user, pass, err := verifyAuth(args.Auth)
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
	if !pass {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusUnauthorized)
		return
	}

	id, err := db.PostMessage(db.Message{
		Sender:  user,
		Content: content,
	})
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}

	msg, err := (&atmt.MessageBuilder[atmt.DataPush]{
		Type: atmt.MessagePush,
		Data: atmt.DataPush{
			MessageID: id,
			Sender:    user,
		},
		Content: content,
	}).Message()
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
	b := atmt.MessageBuilder[atmt.DataResponseBuilder[trinity.ValuesPostMessage]]{
		Type: atmt.MessageResponse,
		Data: atmt.DataResponseBuilder[trinity.ValuesPostMessage]{
			StatusCode: atmt.StatusOK,
			Values: trinity.ValuesPostMessage{
				MessageID: id,
			},
		},
		Content: content,
	}
	if err = b.Write(resp); err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}

	if err = client.PushMessageToListeners(msg); err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
}

func handleGetMessage(resp *atmt.Message, req atmt.DataRequest[trinity.Action]) {
	logPrefix := "handle get message:"
	var args trinity.ArgsGetMessage
	if err := json.Unmarshal(req.Args, &args); err != nil {
		logs.Warning(logPrefix, err)
		atmt.Error(resp, atmt.StatusBadRequest)
		return
	}
	_, pass, err := verifyAuth(args.Auth)
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
	if !pass {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusUnauthorized)
		return
	}

	msg, err := db.GetMessage(args.ID)
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
	b := atmt.MessageBuilder[atmt.DataResponseBuilder[trinity.ValuesGetMessage]]{
		Type: atmt.MessageResponse,
		Data: atmt.DataResponseBuilder[trinity.ValuesGetMessage]{
			StatusCode: atmt.StatusOK,
			Values: trinity.ValuesGetMessage{
				Sender:    msg.Sender,
				MessageID: msg.ID,
			},
		},
		Content: msg.Content,
	}
	if err = b.Write(resp); err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
}

func handleQueryMessageIdsLatestCount(resp *atmt.Message, req atmt.DataRequest[trinity.Action]) {
	logPrefix := "handle query message ids latest count:"
	var args trinity.ArgsQueryMessageIdsLatestCount
	if err := json.Unmarshal(req.Args, &args); err != nil {
		logs.Warning(logPrefix, err)
		atmt.Error(resp, atmt.StatusBadRequest)
		return
	}
	_, pass, err := verifyAuth(args.Auth)
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
	if !pass {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusUnauthorized)
		return
	}

	ids, err := db.QueryMessageIdsLatestCount(args.Count)
	if err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}

	b := atmt.MessageBuilder[atmt.DataResponseBuilder[trinity.ValuesQueryMessageIdsLatestCount]]{
		Type: atmt.MessageResponse,
		Data: atmt.DataResponseBuilder[trinity.ValuesQueryMessageIds]{
			StatusCode: atmt.StatusOK,
			Values: trinity.ValuesQueryMessageIds{
				Ids: ids,
			},
		},
	}
	if *resp, err = b.Message(); err != nil {
		logs.Error(logPrefix, err)
		atmt.Error(resp, atmt.StatusInternalServerError)
		return
	}
}